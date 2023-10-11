package main

import (
	"fmt"
	"github.com/i401/goHive5"
)

func main() {
	handler := thehive5.CreateLogin("https://hive.uauth.io/", "hunter2", true)

	observables := &[]thehive5.Observable{
		{DataType: "ip", Data: "8.8.8.8"},
		{DataType: "domain", Data: "google.com"},
	}

	// Create a new empty customField slice
	customFields := &[]thehive5.CustomField{{Name: "UUID", Group: "Group", Description: "UUID", Type: "string", Value: uuid.New()}}

	alertObject := &thehive5.HiveAlert{
		Title:        "Alert Title",
		Description:  "Alert Description",
		Observables:  observables,
		Status:       "InProgress",
		Tlp:          thehive5.TlpAmber.String(),
		Pap:          thehive5.PapRed.String(),
		Severity:     thehive5.SeverityHigh.String(),
		Tags:         []string{"example", "tag"},
		Source:       "Defender for Endpoint",
		SourceRef:    "#123123124",
		ExternalLink: "https://uauth.io",
		CustomFields: customFields,
		Flag:         true,
	}

	ret, err := handler.CreateAlert(alertObject)
	if err != nil {
		fmt.Printf("%+v", err)
	}

	fmt.Println(ret)
}
