package main

import (
	"fmt"
	"github.com/b401/goHive5"
	"time"
)

func main() {
	handler := thehive5.CreateLogin("https://hive.uauth.io/", "hunter2", true)

	// Simple datatype example
	ipObservable := &thehive5.Observable{
		DataType:         "ip",
		Data:             "127.0.0.1",
		Message:          "Testing",
		Tlp:              "amber+strict", // or thehive5.TlpAmber_Strict.String()
		Pap:              "clear",
		Tags:             []string{"DANGER"},
		Ioc:              true,
		Sighted:          true,
		SightedAt:        time.Now().Add(-24 * time.Hour),
		IgnoreSimilarity: true,
	}

	caseId := 6108
	err := handler.AddCaseObservable(caseId, ipObservable)
	if err != nil {
		fmt.Println(err)
	}
}
