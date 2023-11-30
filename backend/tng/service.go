package tng

import (
	"api-oa-integrator/internal/database"
	"api-oa-integrator/utils"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"maps"
	"net/http"
	"time"
)

type Config struct {
	database.IntegratorConfig
	PlazaId string
}

func (c Config) VerifyVehicle(plateNumber, entryLane string) error {
	zap.L().Sugar().With("plateNumber", plateNumber).Info("VerifyVehicle")
	if plateNumber == "" {
		return errors.New("empty plate number")
	}

	signature := createSignature("ssh/id_rsa.pub")
	if signature == "" {
		return errors.New("empty signature")
	}

	extendInfo, err := json.Marshal(map[string]any{
		"vehiclePlateNo": plateNumber,
		"vehicleType":    "Motorcar",
	})
	reqBody := map[string]interface{}{
		"request": map[string]any{
			"header": map[string]any{
				"requestId": uuid.New().String(),
				"timestamp": time.Now().Format(time.RFC3339),
				"clientId":  c.ClientID.String,
				"function":  "falcon.device.status",
				"version":   viper.GetString("app.version"),
			},
			"body": map[string]any{
				"deviceInfo": map[string]any{
					"deviceType": deviceTypeLPR,
					"deviceNo":   plateNumber,
				},
				"entryTimestamp": time.Now().Format(time.RFC3339),
				"entrySPId":      c.SpID.String,
				"entryPlazaId":   c.PlazaId,
				"entryLaneId":    entryLane,
				"extendInfo":     fmt.Sprintf("%v", string(extendInfo)),
			},
		},
		"signature": signature,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		zap.L().Sugar().With("plateNumber", plateNumber).Errorf("Error marshaling data to JSON: %v", err)
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%v/falcon/device/status", c.Url.String), bytes.NewBuffer(jsonData))
	if err != nil {
		zap.L().Sugar().With("plateNumber", plateNumber).Errorf("Error creating request: %v", err)
		return err
	}

	client := &http.Client{}
	client.Transport = &utils.LoggingRoundTripper{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	var data map[string]any
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return err
	}

	responseBody := data["response"].(map[string]any)["body"]
	if responseBody.(map[string]any)["responseInfo"].(map[string]any)["responseCode"].(string) != "000" {
		return errors.New(fmt.Sprintf("fail to verify vehicle %v", responseBody))
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("invalid response status")
	}
	return nil
}

type TransactionArg struct {
	LPN       string
	EntryLane string
	ExitLane  string
	Amount    float64
	EntryTime time.Time
}

func (c Config) PerformTransaction(locationId, plateNumber, entryLane, exitLane string, entryAt time.Time, amount float64) (map[string]any, map[string]any, error) {
	zap.L().Sugar().With("plateNumber", plateNumber).Info("VerifyVehicle")
	if plateNumber == "" {
		return nil, nil, errors.New("empty plate number")
	}

	signature := createSignature("ssh/id_rsa.pub")
	if signature == "" {
		return nil, nil, errors.New("empty signature")
	}

	extendInfo, err := json.Marshal(map[string]any{
		"vehiclePlateNo": plateNumber,
		"vehicleType":    "Motorcar",
	})

	taxData := map[string]any{
		"surchargeAmt":    0.00,
		"surchargeTaxAmt": 0.00,
		"parkingAmt":      amount, // not sure what is this.
		"parkingTaxAmt":   0.00,   // not sure what is this.
	}
	now := time.Now()
	body := map[string]any{
		"deviceInfo": map[string]any{
			"deviceType": deviceTypeLPR,
			"deviceNo":   plateNumber,
		},
		"serialNum":       fmt.Sprintf("3%v%v%v%v00", c.SpID.String, c.PlazaId, exitLane, now.Format("20060102150405")),
		"transactionType": "C", //Complete (Closed System – populate the Entry and Exit information)
		"entryTimestamp":  entryAt,
		"entrySPId":       c.SpID.String,
		"entryPlazaId":    c.PlazaId,
		"entryLaneId":     entryLane,
		"appSector":       appSectorParking,
		"exitTimestamp":   now.Format(time.RFC3339),
		"exitSPId":        c.SpID.String,
		"exitPlazaId":     c.PlazaId,
		"exitLaneId":      exitLane,
		"vehicleClass":    vehicleClassPrivate,
		"tranAmt":         amount,
		"extendInfo":      fmt.Sprintf("%v", string(extendInfo)),
	}

	maps.Copy(body, taxData)
	reqBody := map[string]interface{}{
		"request": map[string]any{
			"header": map[string]any{
				"requestId": uuid.New().String(),
				"timestamp": time.Now().Format(time.RFC3339),
				"clientId":  c.ClientID.String,
				"function":  "falcon.parking.transaction",
				"version":   viper.GetString("app.version"),
			},
			"body": body,
		},
		"signature": signature,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		zap.L().Sugar().With("plateNumber", plateNumber).Errorf("Error marshaling data to JSON: %v", err)
		return body, taxData, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%v/falcon/parking/transaction", c.Url.String), bytes.NewBuffer(jsonData))
	if err != nil {
		zap.L().Sugar().With("plateNumber", plateNumber).Errorf("Error creating request: %v", err)
		return body, taxData, err
	}

	client := &http.Client{}
	client.Transport = &utils.LoggingRoundTripper{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	resp, err := client.Do(req)
	if err != nil {
		return body, taxData, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return body, taxData, errors.New("invalid response status")
	}
	var data map[string]any
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return body, taxData, err
	}
	responseBody := data["response"].(map[string]any)["body"]
	if responseBody.(map[string]any)["responseInfo"].(map[string]any)["responseCode"].(string) != "000" {
		return body, taxData, errors.New(fmt.Sprintf("fail to perform transaction %v", responseBody))
	}
	return body, taxData, nil
}
