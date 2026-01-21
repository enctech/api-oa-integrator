package tracing

import (
	"go.opentelemetry.io/otel/attribute"
)

const (
	TransactionIDKey = attribute.Key("business.transaction_id")
	JobIDKey         = attribute.Key("business.job_id")
	FacilityKey      = attribute.Key("business.facility")
	DeviceKey        = attribute.Key("business.device")
	PlateNumberKey   = attribute.Key("business.plate_number")
	CustomerIDKey    = attribute.Key("business.customer_id")
	VendorKey        = attribute.Key("business.vendor")
	EntryLaneKey     = attribute.Key("business.entry_lane")
	ExitLaneKey      = attribute.Key("business.exit_lane")
	AmountKey        = attribute.Key("business.amount")
)

func TransactionAttributes(transactionID, plateNumber, vendor string) []attribute.KeyValue {
	attrs := make([]attribute.KeyValue, 0, 3)
	if transactionID != "" {
		attrs = append(attrs, TransactionIDKey.String(transactionID))
	}
	if plateNumber != "" {
		attrs = append(attrs, PlateNumberKey.String(plateNumber))
	}
	if vendor != "" {
		attrs = append(attrs, VendorKey.String(vendor))
	}
	return attrs
}

func JobAttributes(jobID, facility, device, plateNumber string) []attribute.KeyValue {
	attrs := make([]attribute.KeyValue, 0, 4)
	if jobID != "" {
		attrs = append(attrs, JobIDKey.String(jobID))
	}
	if facility != "" {
		attrs = append(attrs, FacilityKey.String(facility))
	}
	if device != "" {
		attrs = append(attrs, DeviceKey.String(device))
	}
	if plateNumber != "" {
		attrs = append(attrs, PlateNumberKey.String(plateNumber))
	}
	return attrs
}
