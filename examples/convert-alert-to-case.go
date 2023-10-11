package main

import (
	"fmt"
	"github.com/b401/goHive5"
)

func main() {
	handler := thehive5.CreateLogin("https://hive.uauth.io", "hunter2", true)

	caseObject := &thehive5.HiveCase{
		Title:       "case title",
		Description: "case description",
		Severity:    "critical",
		Tlp:         "amber",
		Pap:         "amber",
		Flag:        true,
	}

	// AlertIds are always strings
	alertId := "~1155176"

	// returns a HiveCaseResponse
	ret, err := handler.CreateCaseFromAlert(alertId, caseObject)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("v+%", ret)
}
