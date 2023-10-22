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
