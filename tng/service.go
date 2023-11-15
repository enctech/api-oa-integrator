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
	"net/http"
	"time"
)

type Config struct {
	database.IntegratorConfig
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
				"entryPlazaId":   c.PlazaID.String,
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

	if data["response"].(map[string]any)["body"].(map[string]any)["responseInfo"].(map[string]any)["responseCode"].(string) != "000" {
		return errors.New("fail to verify vehicle")
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

func (c Config) PerformTransaction(in TransactionArg) error {
	zap.L().Sugar().With("plateNumber", in.LPN).Info("VerifyVehicle")
	if in.LPN == "" {
		return errors.New("empty plate number")
	}

	signature := createSignature("ssh/id_rsa.pub")
	if signature == "" {
		return errors.New("empty signature")
	}

	extendInfo, err := json.Marshal(map[string]any{
		"vehiclePlateNo": in.LPN,
		"vehicleType":    "Motorcar",
	})
	now := time.Now()
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
					"deviceNo":   in.LPN,
				},
				"serialNum":       fmt.Sprintf("3%v%v%v%v0", c.SpID.String, c.PlazaID.String, in.ExitLane, now.Format("20060102150405")),
				"transactionType": "C", //Complete (Closed System – populate the Entry and Exit information)
				"entryTimestamp":  in.EntryTime,
				"entrySPId":       c.SpID.String,
				"entryPlazaId":    c.PlazaID.String,
				"entryLaneId":     in.EntryLane,
				"appSector":       "09", //Defaults to 09 (Parking)
				"exitTimestamp":   now.Format(time.RFC3339),
				"exitSPId":        c.SpID.String,
				"exitPlazaId":     c.PlazaID.String,
				"exitLaneId":      in.ExitLane,
				"vehicleClass":    "01", //Private Cars (Vehicles with two axles and three or four wheels (excluding taxi and bus))
				"tranAmt":         in.Amount,
				"surchargeAmt":    0.00,
				"surchargeTaxAmt": 0.00,
				"parkingAmt":      in.Amount, // not sure what is this.
				"parkingTaxAmt":   0.00,      // not sure what is this.
				"extendInfo":      extendInfo,
			},
		},
		"signature": signature,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		zap.L().Sugar().With("plateNumber", in.LPN).Errorf("Error marshaling data to JSON: %v", err)
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%v/falcon/parking/transaction", c.Url.String), bytes.NewBuffer(jsonData))
	if err != nil {
		zap.L().Sugar().With("plateNumber", in.LPN).Errorf("Error creating request: %v", err)
		return err
	}

	client := &http.Client{}
	client.Transport = &utils.LoggingRoundTripper{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("invalid response status")
	}
	return nil
}
