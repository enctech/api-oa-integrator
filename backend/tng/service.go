package tng

import (
	"api-oa-integrator/database"
	"api-oa-integrator/logger"
	"api-oa-integrator/utils"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"maps"
	"math"
	"net/http"
	"strconv"
	"time"
)

type Config struct {
	database.IntegratorConfig
	PlazaId string
}

func (c Config) VerifyVehicle(plateNumber, entryLane string) error {
	logger.LogData("info", "VerifyVehicle", map[string]interface{}{"plateNumber": plateNumber, "vendor": "tng"})

	if plateNumber == "" {
		return errors.New("empty plate number")
	}

	var extra map[string]string
	_ = json.Unmarshal(c.Extra.RawMessage, &extra)
	signature := createSignature(extra["sshKey"])
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
		logger.LogData("error", fmt.Sprintf("Error marshaling data to JSON: %v", err), map[string]interface{}{"plateNumber": plateNumber, "vendor": "tng"})
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%v/falcon/device/status", c.Url.String), bytes.NewBuffer(jsonData))
	if err != nil {
		logger.LogData("error", fmt.Sprintf("Error creating request: %v", err), map[string]interface{}{"plateNumber": plateNumber, "vendor": "tng"})
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
	logger.LogData("error", "PerformTransaction", map[string]interface{}{"plateNumber": plateNumber, "vendor": "tng"})
	if plateNumber == "" {
		return nil, nil, errors.New("empty plate number")
	}

	var extra map[string]string
	_ = json.Unmarshal(c.Extra.RawMessage, &extra)
	signature := createSignature(extra["sshKey"])
	if signature == "" {
		return nil, nil, errors.New("empty signature")
	}

	extendInfo, err := json.Marshal(map[string]any{
		"vehiclePlateNo": plateNumber,
		"vehicleType":    "Motorcar",
	})

	tax := c.calculateTax(amount)

	taxData := map[string]any{
		"surcharge":       tax.surcharge,
		"surchargeAmt":    tax.surchargeAmt,
		"surchargeTaxAmt": tax.surchargeTaxAmt,
		"parkingAmt":      tax.parkingAmt,
		"parkingTaxAmt":   tax.parkingTaxAmt,
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
		logger.LogData("error", fmt.Sprintf("Error marshaling data to JSON: %v", err), map[string]interface{}{"plateNumber": plateNumber, "vendor": "tng"})
		return body, taxData, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%v/falcon/parking/transaction", c.Url.String), bytes.NewBuffer(jsonData))
	if err != nil {
		logger.LogData("error", fmt.Sprintf("Error creating request: %v", err), map[string]interface{}{"plateNumber": plateNumber, "vendor": "tng"})
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

type TaxCalculation struct {
	surcharge       float64
	surchargeAmt    float64
	surchargeTaxAmt float64
	parkingAmt      float64
	parkingTaxAmt   float64
}

func (c Config) calculateTax(txn float64) TaxCalculation {
	tax, _ := strconv.ParseFloat(c.TaxRate.String, 64)
	surcharge, _ := strconv.ParseFloat(c.Surcharge.String, 64)
	if c.SurchangeType.SurchargeType == database.SurchargeTypeExact {
		return calculateExactSurchargeAmount(txn, tax, surcharge)
	}

	return calculatePercentSurchargeAmount(txn, tax, surcharge)
}

func calculateExactSurchargeAmount(txn, tax, surcharge float64) TaxCalculation {
	surcF := surcharge
	taxPerc := (tax / 100) * 100
	surchargeAmt := surcF * (100 / (100 + taxPerc))
	surchargeTaxAmt := surcF * (tax / (100 + taxPerc))
	parkingAmt := (txn - surcF) * (100 / (100 + taxPerc))
	parkingTaxAmt := (txn - surcF) * (tax / (100 + taxPerc))
	return TaxCalculation{
		surcharge:       roundMoney(surcF),
		surchargeAmt:    roundMoney(surchargeAmt),
		surchargeTaxAmt: roundMoney(surchargeTaxAmt),
		parkingAmt:      roundMoney(parkingAmt),
		parkingTaxAmt:   roundMoney(parkingTaxAmt),
	}
}

func calculatePercentSurchargeAmount(txn, tax, surcharge float64) TaxCalculation {
	_surcharge := roundMoney((txn - (txn * (surcharge / (100 + tax)))) * (tax / (100 + tax)))

	surchargeAmt := roundMoney((txn * (surcharge / (100 + tax))) * (100 / (100 + tax)))

	surchargeTaxAmt := roundMoney((txn * (surcharge / (100 + tax))) * (surcharge / (100 + tax)))

	parkingAmt := roundMoney((txn - (txn * (surcharge / (100 + tax)))) * (100 / (100 + tax)))

	parkingTaxAmt := roundMoney((txn - (txn * (surcharge / (100 + tax)))) * (tax / (100 + tax)))
	return TaxCalculation{
		surcharge:       _surcharge,
		surchargeAmt:    surchargeAmt,
		surchargeTaxAmt: surchargeTaxAmt,
		parkingAmt:      parkingAmt,
		parkingTaxAmt:   parkingTaxAmt,
	}
}

// Copy from https://stackoverflow.com/a/29786394
func roundMoney(amount float64) float64 {
	return toFixed(amount, 2)
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
