package main

import (
	"fmt"
	"github.com/b401/goHive5"
)

func main() {
	handler := thehive5.CreateLogin("https://hive.uauth.io/", "hunter2", true)
	caseId := 5989

	ret, err := handler.GetCase(caseId)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v\n", ret)
}
