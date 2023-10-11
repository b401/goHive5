package main

import (
	"fmt"
	"github.com/b401/goHive5"
	"os"
	"time"
)

func main() {
	handler := thehive5.CreateLogin("https://hive.uauth.io/", "hunter2", true)

	// open a file
	file, err := os.Open("./malware.zip")
	if err != nil {
		fmt.Println("error!")
		os.Exit(1)
	}
	// Close it at the end
	defer file.Close()

	// Initialize empty Observable
	observable := *new(thehive5.Observable)
	observable.DataType = "file" // or any other dataType that isAttachment = true
	observable.Message = "Some Message"
	// You can use either the TLP/PAP constants or your own string
	observable.Tlp = thehive5.TlpAmber.String()
	observable.Pap = "amber"
	// Define if you wanna unzip the file on theHive after uploading
	observable.IsZip = true
	// Password to unzip the file
	observable.ZipPassword = "malware"
	observable.Tags = []string{"malware"}
	observable.Sighted = true
	observable.SightedAt = time.Now()

	caseId := 6108

	// Only returns data on failure
	err = handler.AddCaseObservableFile(caseId, &observable, file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
