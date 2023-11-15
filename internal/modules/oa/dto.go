package oa

import "encoding/xml"

type VersionRequestWrapper struct {
	Version struct {
		XMLName                xml.Name      `xml:"version" json:"-"`
		EntervoVersion         string        `xml:"entervoVersion"`
		SbAuthorizationVersion string        `xml:"sbAuthorizationVersion"`
		Configuration          Configuration `xml:"configuration"`
	} `xml:"version" json:"version"`
}

type MediaDataWrapper struct {
	MediaDataRequest struct {
		XMLName   xml.Name `xml:"mediaData" json:"-"`
		HashValue struct {
			Value string `xml:"value"`
		} `xml:"hashValue"`
		MediaType string `xml:"mediaType"`
	} `xml:"mediaData" json:"mediaData"`
}

type FinalMessageSBWrapper struct {
	FinalMessageSB FinalMessageSB `xml:"finalMessageSB"`
}

type CancelJobWrapper struct {
	CancelJobWrapper struct {
		XMLName xml.Name `xml:"cancel" json:"-"`
		Reason  struct {
			CancelCode string `xml:"cancelCode"`
			ReasonText string `xml:"reasonText"`
		} `xml:"reason"`
	} `xml:"cancel" json:"cancel"`
}

type FinalMessageSB struct {
	XMLName      xml.Name `xml:"finalMessageSB" json:"-"`
	PaymentMedia string   `xml:"paymentMedia"`
	FinalState   string   `xml:"finalState"`
}

type VersionResponse struct {
	XMLName         xml.Name      `xml:"version"`
	Configuration   Configuration `xml:"configuration"`
	CustomerVersion string        `xml:"customerVersion"`
}

type ConfirmationResponse struct {
	XMLName                  xml.Name `xml:"confirmation"`
	ConfirmationDetailStatus string   `xml:"confirmationDetailStatus"`
	ConfirmationStatus       string   `xml:"confirmationStatus"`
}

type Configuration struct {
	SupportedFunctions []string `xml:"supportedFunctions"`
}

type RequestMetadata struct {
	device   string `param:"device"`
	facility string `param:"facility"`
	jobId    string `param:"jobId"`
}

type JobWrapper struct {
	Job Job `xml:"job"`
}

type Job struct {
	XMLName xml.Name `xml:"job" json:"-"`
	JobType string   `xml:"jobType"`
	JobId   struct {
		ID string `xml:"id"`
	} `xml:"jobId"`
	BusinessTransaction *BusinessTransaction `xml:"businessTransaction"`
	MediaDataList       struct {
		MediaType  string `xml:"mediaType"`
		Identifier struct {
			Name string `xml:"name"`
		} `xml:"identifier"`
	} `xml:"mediaDataList"`
	TimeAndPlace struct {
		Operator struct {
			OperatorNumber string `xml:"operatorNumber"`
		} `xml:"operator"`
		Computer struct {
			ComputerNumber string `xml:"computerNumber"`
		} `xml:"computer"`
		Device struct {
			DeviceNumber string `xml:"deviceNumber"`
			DeviceType   string `xml:"deviceType"`
		} `xml:"device"`
		TransactionTimeStamp struct {
			TimeStamp string `xml:"timeStamp"`
		} `xml:"transactionTimeStamp"`
		Facility struct {
			FacilityNumber string `xml:"facilityNumber"`
		} `xml:"facility"`
	} `xml:"timeAndPlace"`
	ProviderInformation struct {
		Provider struct {
			ProviderId   string `xml:"providerId"`
			ProviderName string `xml:"providerName"`
		} `xml:"provider"`
	} `xml:"providerInformation"`
	CustomerInformation *CustomerInformation `xml:"customerInformation"`
	PaymentData         *PaymentData         `xml:"paymentData"`
}

type (
	OriginalAmount struct {
		Amount  string `xml:"amount"`
		VatRate string `xml:"vatRate"`
	}
	PaymentData struct {
		OriginalAmount  OriginalAmount `xml:"originalAmount"`
		RemainingAmount struct {
			Text    string `xml:",chardata"`
			Amount  string `xml:"amount"`
			VatRate string `xml:"vatRate"`
		} `xml:"remainingAmount"`
	}
	PaymentInformation struct {
		PaymentLocation string       `xml:"paymentLocation"`
		PayedAmount     *PayedAmount `xml:"payedAmount"`
	}

	PayedAmount struct {
		Amount  string `xml:"amount"`
		VatRate string `xml:"vatRate"`
	}

	Provider struct {
		ProviderId   string `xml:"providerId"`
		ProviderName string `xml:"providerName"`
	}

	ProviderInformation struct {
		Provider Provider `xml:"provider"`
	}

	Customer struct {
		CustomerId    string `xml:"customerId"`
		CustomerName  string `xml:"customerName"`
		CustomerGroup string `xml:"customerGroup"`
	}

	CustomerInformation struct {
		Customer Customer `xml:"customer"`
	}

	ReservationTariff struct {
		TariffName   string `xml:"tariffName"`
		TariffNumber int    `xml:"tariffNumber"`
	}

	Reservation struct {
		ReservationTariff ReservationTariff `xml:"reservationTariff"`
	}

	Identifier struct {
		Name string `xml:"name"`
	}

	MediaDataList struct {
		MediaType  string     `xml:"mediaType"`
		Identifier Identifier `xml:"identifier"`
	}

	BusinessTransaction struct {
		ID string `xml:"id"`
	}

	FinalMessageCustomer struct {
		XMLName             xml.Name             `xml:"finalMessageCustomer"`
		PaymentInformation  *PaymentInformation  `xml:"paymentInformation"`
		ProviderInformation *ProviderInformation `xml:"providerInformation"`
		CustomerInformation *CustomerInformation `xml:"customerInformation"`
		MediaDataList       *[]MediaDataList     `xml:"mediaDataList"`
		Counting            *string              `xml:"counting"`
		BusinessTransaction *BusinessTransaction `xml:"businessTransaction"`
	}
)
