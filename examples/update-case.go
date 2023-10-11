package main

import (
	"fmt"
	"github.com/b401/goHive5"
)

func main() {
	handler := thehive5.CreateLogin("https://hive.uauth.io/", "hunter2", true)

	flag := false
	caseId := 6110

	caseObject := &thehive5.HiveUpdateCase{
		Title:        "case title",
		Description:  "case description",
		Tlp:          "clear",
		AddTags:      []string{"newTag"},
		RemoveTags:   []string{"oldTag"},
		Tags:         []string{"gohive5", "example", "case"},
		ImpactStatus: "withImpact",
		Flag:         &flag,
	}

	err := handler.UpdateCase(caseId, caseObject)
	if err != nil {
		// do some error handling
		fmt.Println(err)
	}
}
