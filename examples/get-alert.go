package main

import (
	"fmt"
	"github.com/b401/goHive5"
)

func main() {
	alertId := "~496406592"
	handler := thehive5.CreateLogin("https://hive.uauth.io/", "hunter2", true)

	// returns a HiveHiveHiveAlertResponse
	ret, err := handler.GetAlert(alertId)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v\n", ret)
}
