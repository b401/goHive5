# goHive5

goHive5 is an unofficial API client library for the [TheHive5](https://www.strangebee.com/thehive/) platform.

Some functionalities are still missing.   
Please feel free to contribute or create issues for additional functionality requests.

Overview: [goHive5](./examples/hive5go_overview.md)

Check out [examples](./examples/) for more examples.


### TODO
* Write tests for functions
* Add more examples
* Add administrative functions
* Add more API coverage

## Install
```Bash
go get github.com/b401/gohive5
```

```Go
import "github.com/b401/gohive5"
```


## Initialize a hive object

```Go
verifyCert := true
hive := thehive5.CreateLogin("https://thehive.example.com", "apitoken", verifyCert) 
```

## Create case example with customFields
```Go
tasks := []thehive5.CaseTask{
	thehive5.CaseTask{Title: "Identification", Status: "Waiting",Flag: true},
	thehive5.CaseTask{Title: "Containment", Description: "Please contain this threat"},
	thehive5.CaseTask{Title: "Eradication", Status: "InProgress", Mandatory: true},
}

caseObject := &thehive5.HiveCase { 
	Title: "case title",
	Description: "case description",
	Severity: "critical",
	Tlp: "amber",
	Pap: "amber",
	Tasks: &tasks,
	Tags: []string{"gohive5", "example", "case"},
	Flag: true,
}

ret, err := handler.CreateCase(caseObject); if err != nil {
	fmt.Println(err)
}
```

This will return a pointer to a HiveCaseResponse struct

```Go
type HiveCaseResponse struct {
	Id                  string        `json:"_id"`
	Title               string        `json:"title"`
	Number              int           `json:"number"`
	Description         string        `json:"description"`
	Status              string        `json:"status"`
	Stage               string        `json:"stage"`
	StartDate           time.Time     `json:"startDate"`
	Tlp                 int           `json:"tlp"`
	Pap                 int           `json:"pap"`
	Type                string        `json:"_type"`
	CreatedBy           string        `json:"_createdBy"`
	UpdatedBy           string        `json:"_updatedBy"`
	CreatedAt           time.Time     `json:"_createdAt"`
	UpdatedAt           time.Time     `json:"_updatedAt"`
	EndDate             time.Time     `json:"endDate"`
	Tags                []string      `json:"tags"`
	Flag                bool          `json:"flag"`
	TlpLabel            string        `json:"tlpLabel"`
	PapLabel            string        `json:"papLabel"`
	Summary             string        `json:"summary"`
	Severity            int           `json:"severity"`
	ImpactStatus        string        `json:"impactStatus"`
	Assignee            string        `json:"assignee"`
	CustomFields        []CustomField `json:"customFields"`
	UserPermissions     []string      `json:"userPermissions"`
	ExtraData           map[string]string 	  `json:"extraData"`
	NewDate             time.Time         `json:"newDate"`
	InProgressDate      time.Time         `json:"inProgressDate"`
	ClosedDate          time.Time         `json:"closedDate"`
	AlertDate           time.Time         `json:"alertDate"`
	AlertNewDate        time.Time         `json:"alertNewDate"`
	AlertInProgressDate time.Time         `json:"alertInProgressDate"`
	AlertImportedDate   time.Time         `json:"alertImportedDate"`
	TimeToDetect        time.Duration         `json:"timeToDetect"`
	TimeToTriage        time.Duration         `json:"timeToTriage"`
	TimeToQualify       time.Duration         `json:"timeToQualify"`
	TimeToAcknowledge   time.Duration         `json:"timeToAcknowledge"`
	TimeToResolve       time.Duration         `json:"timeToResolve"`
	HandlingDuration    time.Duration         `json:"handlingDuration"`
}
```

## Create alert example with customFields & observables

```Go
observables := 	&[]thehive5.Observable {
	thehive5.Observable{DataType: "ip", Data: "8.8.8.8"},
	thehive5.Observable{DataType: "domain", Data: "google.com"},
}

// Create a new empty customField slice
customFields := &[]thehive5.CustomField{ thehive5.CustomField{ Name: "UUID", Group: "Group", Description: "UUID", Type: "string", Value: uuid.New()}}

alertObject := &thehive5.HiveAlert {
	Title: "Alert Title",
	Description: "Alert Description",
	Observables: observables,
	Status: "InProgress",
	Tlp: thehive5.TlpAmber.String(),
	Pap: thehive5.PapRed.String(),
	Severity: thehive5.SeverityHigh.String(),
	Tags: []string{"example", "tag"},
	Source: "Defender for Endpoint",
	SourceRef: "#123123124",
	ExternalLink: "https://uauth.io",
	CustomFields: customFields,
	Flag: true,
}

ret, err := handler.CreateAlert(alertObject); if err != nil {
	fmt.Printf("%+v", err)
}

```

``hive.CreateAlert`` returns a pointer to an HiveAlertResponse struct



