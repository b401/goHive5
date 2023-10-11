package main

import (
	"fmt"
	"github.com/b401/goHive5"
)

func main() {
	handler := thehive5.CreateLogin("https://hive.uauth.io/", "hunter2", true)

	tasks := []thehive5.CaseTask{
		{Title: "Identification", Status: "Waiting", Flag: true},
		{Title: "Containment", Description: "Please contain this threat"},
		{Title: "Eradication", Status: "InProgress", Mandatory: true},
	}

	caseObject := &thehive5.HiveCase{
		Title:       "case title",
		Description: "case description",
		Severity:    "critical",
		Tlp:         "amber",
		Pap:         "amber",
		Tasks:       &tasks,
		Tags:        []string{"gohive5", "example", "case"},
		Flag:        true,
	}

	// ret contains now a HiveCaseResponse struct
	// We populate all fields with standard values (Mean times are in time.Duration & every time object is time.Time - use .IsZero())
	ret, err := handler.CreateCase(caseObject)
	if err != nil {
		// do some error handling
		fmt.Println(err)
	}

	fmt.Println(ret)

}
