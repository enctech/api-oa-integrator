package tng

import (
	"api-oa-integrator/database"
	"api-oa-integrator/logger"
	"api-oa-integrator/tracing"
	"api-oa-integrator/utils"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Config struct {
	database.IntegratorConfig
	PlazaId string
}

const statusCodeSuccess = "000"
const statusCodeDuplicateTransaction = "998"
const timeOut = time.Second * 5
const transactionTimeOut = time.Second * 3
const voidDelayDuration = time.Second * 5

func (c Config) VerifyVehicle(ctx context.Context, plateNumber, entryLane string) error {
	ctx, span := tracing.Tracer().Start(ctx, "tng.VerifyVehicle",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			tracing.VendorKey.String("tng"),
			tracing.PlateNumberKey.String(plateNumber),
			tracing.EntryLaneKey.String(entryLane),
		),
	)
	defer span.End()

	logger.LogWithContext(ctx, "info", "VerifyVehicle", map[string]interface{}{"plateNumber": plateNumber, "vendor": "tng"})

	if plateNumber == "" {
		err := errors.New("tng error: empty plate number")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	var extra map[string]string
	_ = json.Unmarshal(c.Extra.RawMessage, &extra)

	extendInfo, err := json.Marshal(map[string]any{
		"vehiclePlateNo": plateNumber,
		"vehicleType":    "Motorcar",
	})

	logger.LogWithContext(ctx, "info", fmt.Sprintf("current time: %v", time.Now().Local().Format(time.RFC3339)), nil)
	reqBody := map[string]interface{}{
		"header": map[string]any{
			"requestId": uuid.New().String(),
			"timestamp": time.Now().Local().Format(time.RFC3339),
			"clientId":  c.ClientID.String,
			"function":  "falcon.device.status",
			"version":   viper.GetString("app.version"),
		},
		"body": map[string]any{
			"deviceInfo": map[string]any{
				"deviceType": deviceTypeLPR,
				"deviceNo":   plateNumber,
			},
			"entryTimestamp": time.Now().Local().Format(time.RFC3339),
			"entrySPId":      c.SpID.String,
			"entryPlazaId":   c.PlazaId,
			"entryLaneId":    entryLane,
			"extendInfo":     fmt.Sprintf("%v", string(extendInfo)),
		},
	}

	originalRequestData, err := json.Marshal(reqBody)
	signer, err := NewSigner(extra["sshKey"])
	if err != nil {
		err = fmt.Errorf("error creating signer: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	dataToSign := string(originalRequestData)
	signature, err := signer.Sign(dataToSign)
	if err != nil {
		fmt.Printf("Error signing data: %v\n", err)
		err = fmt.Errorf("error signing data: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	if signature == "" {
		err := errors.New("tng error: empty signature")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	finalRequest := map[string]interface{}{
		"request":   reqBody,
		"signature": signature,
	}
	jsonData, err := json.Marshal(finalRequest)
	if err != nil {
		logger.LogWithContext(ctx, "error", fmt.Sprintf("Error marshaling data to JSON: %v", err), map[string]interface{}{"plateNumber": plateNumber, "vendor": "tng"})
		err = fmt.Errorf("tng error: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	reqCtx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()
	req, err := http.NewRequestWithContext(reqCtx, "POST", fmt.Sprintf("%v/falcon/device/status", c.Url.String), bytes.NewBuffer(jsonData))
	if err != nil {
		logger.LogWithContext(ctx, "error", fmt.Sprintf("Error creating request: %v", err), map[string]interface{}{"plateNumber": plateNumber, "vendor": "tng"})
		err = fmt.Errorf("tng error: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	resp, err := utils.GlobalInsecureHttpClient.Do(req)
	if err != nil {
		err = fmt.Errorf("tng error: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	defer resp.Body.Close()
	var data map[string]any
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		err = fmt.Errorf("tng error: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	responseBody := data["response"].(map[string]any)["body"]
	if responseBody.(map[string]any)["responseInfo"].(map[string]any)["responseCode"].(string) != "000" {
		err = fmt.Errorf("tng error: fail to verify vehicle %v", responseBody)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	if resp.StatusCode != http.StatusOK {
		err := errors.New("tng error: invalid response status")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
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

func (c Config) VoidTransaction(ctx context.Context, plateNumber, transactionId string) error {
	ctx, span := tracing.Tracer().Start(ctx, "tng.VoidTransaction",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			tracing.VendorKey.String("tng"),
			tracing.PlateNumberKey.String(plateNumber),
			tracing.TransactionIDKey.String(transactionId),
		),
	)
	defer span.End()

	logger.LogWithContext(ctx, "info", "VoidTransaction", map[string]interface{}{"plateNumber": plateNumber, "vendor": "tng"})

	body := map[string]any{
		"deviceInfo": map[string]any{
			"deviceType": deviceTypeLPR,
			"deviceNo":   plateNumber,
		},
		"serialNum":         transactionId,
		"cancelRequestTime": time.Now().Local().Format(time.RFC3339),
		"extendInfo":        nil,
	}

	var extra map[string]string
	_ = json.Unmarshal(c.Extra.RawMessage, &extra)

	reqBody := map[string]interface{}{
		"header": map[string]any{
			"requestId": uuid.New().String(),
			"timestamp": time.Now().Local().Format(time.RFC3339),
			"clientId":  c.ClientID.String,
			"function":  "falcon.parking.cancel.transaction.order",
			"version":   viper.GetString("app.version"),
		},
		"body": body,
	}

	originalDataBytes, err := json.Marshal(reqBody)
	signer, _ := NewSigner(extra["sshKey"])
	signature, err := signer.Sign(string(originalDataBytes))
	if signature == "" {
		err := errors.New("tng error: empty signature")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	finalRequest := map[string]interface{}{
		"request":   reqBody,
		"signature": signature,
	}

	jsonData, err := json.Marshal(finalRequest)
	if err != nil {
		logger.LogWithContext(ctx, "error", fmt.Sprintf("Error marshaling data to JSON: %v", err), map[string]interface{}{"plateNumber": plateNumber, "vendor": "tng"})
		err = fmt.Errorf("tng error: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	reqCtx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, "POST", fmt.Sprintf("%v/falcon/parking/cancel/transaction-order", c.Url.String), bytes.NewBuffer(jsonData))
	if err != nil {
		logger.LogWithContext(ctx, "error", fmt.Sprintf("Error creating request: %v", err), map[string]interface{}{"plateNumber": plateNumber, "vendor": "tng"})
		err = fmt.Errorf("tng error: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	resp, err := utils.GlobalInsecureHttpClient.Do(req)
	if err != nil {
		err = fmt.Errorf("tng error: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("tng error: %v", "invalid response status")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	var data map[string]any
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		err = fmt.Errorf("tng error: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	responseBody := data["response"].(map[string]any)["body"]
	statusCode := responseBody.(map[string]any)["responseInfo"].(map[string]any)["responseCode"].(string)
	if statusCode != statusCodeSuccess && statusCode != statusCodeDuplicateTransaction {
		err := fmt.Errorf("fail to perform transaction %v", responseBody)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}

func (c Config) PerformTransaction(ctx context.Context, locationId, plateNumber, entryLane, exitLane string, entryAt time.Time, amount float64) (map[string]any, map[string]any, *string, error) {
	ctx, span := tracing.Tracer().Start(ctx, "tng.PerformTransaction",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			tracing.VendorKey.String("tng"),
			tracing.PlateNumberKey.String(plateNumber),
			tracing.EntryLaneKey.String(entryLane),
			tracing.ExitLaneKey.String(exitLane),
			tracing.AmountKey.Float64(amount),
		),
	)
	defer span.End()

	if plateNumber == "" {
		err := errors.New("tng error: empty plate number")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, nil, nil, err
	}

	var extra map[string]string
	_ = json.Unmarshal(c.Extra.RawMessage, &extra)

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
	serialNum := fmt.Sprintf("3%v%v%v%v00", c.SpID.String, c.PlazaId, exitLane, now.Format("20060102150405"))
	span.SetAttributes(attribute.String("tng.serial_num", serialNum))

	body := map[string]any{
		"deviceInfo": map[string]any{
			"deviceType": deviceTypeLPR,
			"deviceNo":   plateNumber,
		},
		"serialNum":       serialNum,
		"transactionType": "C",
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
		"header": map[string]any{
			"requestId": uuid.New().String(),
			"timestamp": time.Now().Local().Format(time.RFC3339),
			"clientId":  c.ClientID.String,
			"function":  "falcon.parking.transaction",
			"version":   viper.GetString("app.version"),
		},
		"body": body,
	}

	originalDataBytes, err := json.Marshal(reqBody)
	signer, _ := NewSigner(extra["sshKey"])
	signature, err := signer.Sign(string(originalDataBytes))
	if signature == "" {
		err := errors.New("tng error: empty signature")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, nil, nil, err
	}

	finalRequest := map[string]interface{}{
		"request":   reqBody,
		"signature": signature,
	}

	jsonData, err := json.Marshal(finalRequest)
	if err != nil {
		logger.LogWithContext(ctx, "error", fmt.Sprintf("Error marshaling data to JSON: %v", err), map[string]interface{}{"plateNumber": plateNumber, "vendor": "tng"})
		err = fmt.Errorf("tng error: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return body, taxData, nil, err
	}

	reqCtx, cancel := context.WithTimeout(ctx, transactionTimeOut)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, "POST", fmt.Sprintf("%v/falcon/parking/transaction", c.Url.String), bytes.NewBuffer(jsonData))
	if err != nil {
		logger.LogWithContext(ctx, "error", fmt.Sprintf("Error creating request: %v", err), map[string]interface{}{"plateNumber": plateNumber, "vendor": "tng"})
		err = fmt.Errorf("tng error: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return body, taxData, nil, err
	}

	resp, err := utils.GlobalInsecureHttpClient.Do(req)
	if err != nil {
		time.Sleep(voidDelayDuration)
		if voidErr := c.VoidTransaction(ctx, plateNumber, serialNum); voidErr != nil {
			err = fmt.Errorf("fail to void transaction %v", voidErr)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return body, taxData, nil, err
		}
		err = fmt.Errorf("tng error: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return body, taxData, nil, err
	}
	defer resp.Body.Close()
	var data map[string]any
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		err = fmt.Errorf("tng error: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return body, taxData, nil, err
	}
	responseBody := data["response"].(map[string]any)["body"]
	statusCode := responseBody.(map[string]any)["responseInfo"].(map[string]any)["responseCode"].(string)
	if statusCode != statusCodeSuccess && statusCode != statusCodeDuplicateTransaction {
		err := fmt.Errorf("response code is not %v or %v", statusCodeSuccess, statusCodeDuplicateTransaction)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return body, taxData, nil, err
	}
	if statusCode == statusCodeDuplicateTransaction {
		status := "duplicate"
		span.SetAttributes(attribute.String("tng.status", status))
		return body, taxData, &status, nil
	}
	return body, taxData, nil, nil
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
func (c Config) CancelEntry() error {
	return nil
}

func calculateExactSurchargeAmount(txnAmt, tax, surcharge float64) TaxCalculation {
	surcF := surcharge
	taxPerc := (tax / 100) * 100
	surchargeAmt := surcF * (100 / (100 + taxPerc))
	surchargeTaxAmt := surcF * (tax / (100 + taxPerc))
	parkingAmt := (txnAmt - surcF) * (100 / (100 + taxPerc))
	parkingTaxAmt := (txnAmt - surcF) * (tax / (100 + taxPerc))
	return TaxCalculation{
		surcharge:       utils.RoundMoney(surcF),
		surchargeAmt:    utils.RoundMoney(surchargeAmt),
		surchargeTaxAmt: utils.RoundMoney(surchargeTaxAmt),
		parkingAmt:      utils.RoundMoney(parkingAmt),
		parkingTaxAmt:   utils.RoundMoney(parkingTaxAmt),
	}
}

// Formula:
//
//	surcharge=(TranAmt-(TranAmt*(Surc%/(100+Tax%))))*(Tax%/(100+Tax%))
//
//	surchargeAmt=(TranAmt*(Surc%/(100+Tax%)))*(100/(100+Tax%))
//
//	surchargeTaxAmt=(TranAmt*(Surc%/(100+Tax%)))*(Surc%/(100+Tax%))
//
//	parkingAmt=(TranAmt-(TranAmt*(Surc%/(100+Tax%))))*(100/(100+Tax%))
//
//	parkingTaxAmt=(TranAmt-(TranAmt*(Surc%/(100+Tax%))))*(Tax%/(100+Tax%))
func calculatePercentSurchargeAmount(txnAmt, tax, surc float64) TaxCalculation {
	_surcharge := utils.RoundMoney((txnAmt - (txnAmt * (surc / (100 + tax)))) * (tax / (100 + tax)))

	surchargeAmt := utils.RoundMoney((txnAmt * (surc / (100 + tax))) * (100 / (100 + tax)))

	surchargeTaxAmt := utils.RoundMoney((txnAmt * (surc / (100 + tax))) * (surc / (100 + tax)))

	parkingAmt := utils.RoundMoney((txnAmt - (txnAmt * (surc / (100 + tax)))) * (100 / (100 + tax)))

	parkingTaxAmt := utils.RoundMoney((txnAmt - (txnAmt * (surc / (100 + tax)))) * (tax / (100 + tax)))
	return TaxCalculation{
		surcharge:       _surcharge,
		surchargeAmt:    surchargeAmt,
		surchargeTaxAmt: surchargeTaxAmt,
		parkingAmt:      parkingAmt,
		parkingTaxAmt:   parkingTaxAmt,
	}
}
