package main

import (
	"fmt"
	"github.com/b401/gohive5"
)

func main() {
	handler := thehive5.CreateLogin("https://hive.uauth.io", "hunter2", true)

	// AlertIds are always strings
	alertId := "~1155176"
	// CaseIds are always ints
	caseId := 6108

	err := handler.MergeAlert(alertId, caseId)
	if err != nil {
		fmt.Println(err)
	}
}
