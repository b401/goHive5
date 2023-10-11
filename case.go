/*
thehive5 implements functionality to interact with the most recent version of thehive.
https://www.strangebee.com/thehive/
*/
package thehive5

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// A HiveCase includes all informations related to a case
// It also includes the fields for an updated task
type HiveCase struct {
	Title             string              `json:"title"`
	Description       string              `json:"description"`
	Severity          string              `json:"severity,omitempty"`
	StartDate         time.Time           `json:"startDate,omitempty"`
	EndDate           time.Time           `json:"endDate,omitempty"`
	Tags              []string            `json:"tags,omitempty"`
	Flag              bool                `json:"flag,omitempty"`
	Tlp               string              `json:"tlp,omitempty"`
	Pap               string              `json:"pap,omitempty"`
	Status            string              `json:"status,omitempty"`
	Summary           string              `json:"summary,omitempty"`
	Assignee          string              `json:"assignee,omitempty"`
	CustomFields      *[]CustomField      `json:"customFields,omitempty"`
	Template          string              `json:"caseTemplate,omitempty"`
	Pages             *[]Pages            `json:"pages,omitempty"`
	Tasks             *[]CaseTask         `json:"tasks,omitempty"`
	SharingParameters *[]SharingParameter `json:"sharingParameters,omitempty"`
	TaskRule          string              `json:"taskRule,omitempty"`
	ObservableRule    string              `json:"observableRule,omitempty"`
}

// Marshalling the case requests converting the time objects into Unixmilli int64 values
func (hc *HiveCase) MarshalJSON() ([]byte, error) {
	type Alias HiveCase
	var (
		startDateInt64 int64
		endDateInt64   int64
		tlpInt         *int
		papInt         *int
		severityInt    *int
	)
	if !hc.StartDate.IsZero() {
		// We ensure that all data sent to the hive is in UTC format
		startDateInt64 = hc.StartDate.UTC().UnixMilli()
	}
	if !hc.EndDate.IsZero() {
		// We ensure that all data sent to the hive is in UTC format
		endDateInt64 = hc.EndDate.UTC().UnixMilli()
	}

	if len(hc.Severity) != 0 {
		var sev Severity
		err := sev.FromString(hc.Severity)
		if err != nil {
			return nil, err
		}
		tmp := int(sev)
		severityInt = &tmp
	}

	if len(hc.Tlp) != 0 {
		var tlp Tlp
		err := tlp.FromString(hc.Tlp)
		if err != nil {
			return nil, err
		}
		tmp := int(tlp)
		tlpInt = &tmp
	}
	if len(hc.Pap) != 0 {
		var pap Pap
		err := pap.FromString(hc.Pap)
		if err != nil {
			return nil, err
		}
		tmp := int(pap)
		papInt = &tmp
	}

	return json.Marshal(&struct {
		StartDate int64 `json:"startDate,omitempty"`
		EndDate   int64 `json:"endDate,omitempty"`
		Tlp       *int  `json:"tlp,omitempty"`
		Pap       *int  `json:"pap,omitempty"`
		Severity  *int  `json:"severity,omitempty"`
		*Alias
	}{
		StartDate: startDateInt64,
		EndDate:   endDateInt64,
		Tlp:       tlpInt,
		Pap:       papInt,
		Severity:  severityInt,
		Alias:     (*Alias)(hc),
	})
}

type Pages struct {
	Title    string  `json:"title"`
	Content  string  `json:"content"`
	Order    *int    `json:"order,omitempty"`
	Category *string `json:"category"`
}

type HiveUpdateCase struct {
	Title             string              `json:"title,omitempty"`
	Description       string              `json:"description,omitempty"`
	Severity          string              `json:"severity,omitempty"`
	StartDate         time.Time           `json:"startDate,omitempty"`
	EndDate           time.Time           `json:"endDate,omitempty"`
	Tags              []string            `json:"tags,omitempty"`
	Flag              *bool               `json:"flag,omitempty"`
	Tlp               string              `json:"tlp,omitempty"`
	Pap               string              `json:"pap,omitempty"`
	Status            string              `json:"status,omitempty"`
	Summary           string              `json:"summary,omitempty"`
	Assignee          string              `json:"assignee,omitempty"`
	CustomFields      *[]CustomField      `json:"customFields,omitempty"`
	Template          string              `json:"caseTemplate,omitempty"`
	Tasks             *[]CaseTask         `json:"tasks,omitempty"`
	SharingParameters *[]SharingParameter `json:"sharingParameters,omitempty"`
	TaskRule          string              `json:"taskRule,omitempty"`
	ObservableRule    string              `json:"observableRule,omitempty"`
	ImpactStatus      string              `json:"impactStatus,omitempty"`
	AddTags           []string            `json:"addTags,omitempty"`
	RemoveTags        []string            `json:"removeTags,omitempty"`
}

// Marshalling the case requests converting the time objects into Unixmilli int64 values
func (hu *HiveUpdateCase) MarshalJSON() ([]byte, error) {
	type Alias HiveUpdateCase
	var (
		startDateInt64 int64
		endDateInt64   int64
		tlpInt         *int
		severityInt    *int
		papInt         *int
	)

	if !hu.StartDate.IsZero() {
		// We ensure that all data sent to the hive is in UTC format
		startDateInt64 = hu.StartDate.UTC().UnixMilli()
	}
	if !hu.EndDate.IsZero() {
		// We ensure that all data sent to the hive is in UTC format
		endDateInt64 = hu.EndDate.UTC().UnixMilli()
	}

	if len(hu.Severity) != 0 {
		var sev Severity
		err := sev.FromString(hu.Severity)
		if err != nil {
			return nil, err
		}
		tmp := int(sev)
		severityInt = &tmp
	}

	if len(hu.Tlp) != 0 {
		var tlp Tlp
		err := tlp.FromString(hu.Tlp)
		if err != nil {
			return nil, err
		}
		tmp := int(tlp)
		tlpInt = &tmp
	}
	if len(hu.Pap) != 0 {
		var pap Pap
		err := pap.FromString(hu.Pap)
		if err != nil {
			return nil, err
		}
		tmp := int(pap)
		papInt = &tmp
	}

	if len(hu.ImpactStatus) != 0 {
		switch strings.ToLower(hu.ImpactStatus) {
		case "withimpact":
			hu.ImpactStatus = "WithImpact"
		case "noimpact":
			hu.ImpactStatus = "NoImpact"
		default:
			return nil, fmt.Errorf("unknown impact value: %s", hu.ImpactStatus)
		}
	}

	return json.Marshal(&struct {
		StartDate int64 `json:"startDate,omitempty"`
		EndDate   int64 `json:"endDate,omitempty"`
		Tlp       *int  `json:"tlp,omitempty"`
		Pap       *int  `json:"pap,omitempty"`
		Severity  *int  `json:"severity,omitempty"`
		*Alias
	}{
		StartDate: startDateInt64,
		EndDate:   endDateInt64,
		Tlp:       tlpInt,
		Severity:  severityInt,
		Pap:       papInt,

		Alias: (*Alias)(hu),
	})
}

// SharingParameter includes the informations necessary to share a case / alert across organizations
type SharingParameter struct {
	Organisation   string `json:"organisation"`
	Share          *bool  `json:"share,omitempty"`
	Profile        string `json:"profile,omitempty"`
	TaskRule       string `json:"taskRule,omitempty"`
	ObservableRule string `json:"observableRule,omitempty"`
}

// HiveCaseResponse stores the response of a case from thehive
type HiveCaseResponse struct {
	Id                  string            `json:"_id"`
	Title               string            `json:"title"`
	Number              int               `json:"number"`
	Description         string            `json:"description"`
	Status              string            `json:"status"`
	Stage               string            `json:"stage"`
	StartDate           time.Time         `json:"startDate"`
	Tlp                 int               `json:"tlp"`
	Pap                 int               `json:"pap"`
	Type                string            `json:"_type"`
	CreatedBy           string            `json:"_createdBy"`
	UpdatedBy           string            `json:"_updatedBy"`
	CreatedAt           time.Time         `json:"_createdAt"`
	UpdatedAt           time.Time         `json:"_updatedAt"`
	EndDate             time.Time         `json:"endDate"`
	Tags                []string          `json:"tags"`
	Flag                bool              `json:"flag"`
	TlpLabel            string            `json:"tlpLabel"`
	PapLabel            string            `json:"papLabel"`
	Summary             string            `json:"summary"`
	Severity            int               `json:"severity"`
	ImpactStatus        string            `json:"impactStatus"`
	Assignee            string            `json:"assignee"`
	CustomFields        []CustomField     `json:"customFields"`
	UserPermissions     []string          `json:"userPermissions"`
	ExtraData           map[string]string `json:"extraData"`
	NewDate             time.Time         `json:"newDate"`
	InProgressDate      time.Time         `json:"inProgressDate"`
	ClosedDate          time.Time         `json:"closedDate"`
	AlertDate           time.Time         `json:"alertDate"`
	AlertNewDate        time.Time         `json:"alertNewDate"`
	AlertInProgressDate time.Time         `json:"alertInProgressDate"`
	AlertImportedDate   time.Time         `json:"alertImportedDate"`
	TimeToDetect        time.Duration     `json:"timeToDetect"`
	TimeToTriage        time.Duration     `json:"timeToTriage"`
	TimeToQualify       time.Duration     `json:"timeToQualify"`
	TimeToAcknowledge   time.Duration     `json:"timeToAcknowledge"`
	TimeToResolve       time.Duration     `json:"timeToResolve"`
	HandlingDuration    time.Duration     `json:"handlingDuration"`
}

// shadowHiveCaseResponse is used to Unmarshal repsonses into time.Time
type shadowHiveCaseResponse struct {
	Id                  string            `json:"_id"`
	Type                string            `json:"_type"`
	CreatedBy           string            `json:"_createdBy"`
	UpdatedBy           string            `json:"_updatedBy"`
	CreatedAt           int64             `json:"_createdAt"`
	UpdatedAt           int64             `json:"_updatedAt"`
	Number              int               `json:"number"`
	Title               string            `json:"title"`
	Description         string            `json:"description"`
	StartDate           int64             `json:"startDate"`
	EndDate             int64             `json:"endDate"`
	Tags                []string          `json:"tags"`
	Flag                bool              `json:"flag"`
	Tlp                 int               `json:"tlp"`
	Pap                 int               `json:"pap"`
	TlpLabel            string            `json:"tlpLabel"`
	PapLabel            string            `json:"papLabel"`
	Status              string            `json:"status"`
	Stage               string            `json:"stage"`
	Summary             string            `json:"summary"`
	Severity            int               `json:"severity"`
	ImpactStatus        string            `json:"impactStatus"`
	Assignee            string            `json:"assignee"`
	CustomFields        []CustomField     `json:"customFields"`
	UserPermissions     []string          `json:"userPermissions"`
	ExtraData           map[string]string `json:"extraData"`
	NewDate             int64             `json:"newDate"`
	InProgressDate      int64             `json:"inProgressDate"`
	ClosedDate          int64             `json:"closedDate"`
	AlertDate           int64             `json:"alertDate"`
	AlertNewDate        int64             `json:"alertNewDate"`
	AlertInProgressDate int64             `json:"alertInProgressDate"`
	AlertImportedDate   int64             `json:"alertImportedDate"`
	TimeToDetect        int64             `json:"timeToDetect"`
	TimeToTriage        int64             `json:"timeToTriage"`
	TimeToQualify       int64             `json:"timeToQualify"`
	TimeToAcknowledge   int64             `json:"timeToAcknowledge"`
	TimeToResolve       int64             `json:"timeToResolve"`
	HandlingDuration    int64             `json:"handlingDuration"`
}

func (hc *HiveCaseResponse) UnmarshalJSON(data []byte) error {
	shadow := new(shadowHiveCaseResponse)
	err := json.Unmarshal(data, &shadow)
	if err != nil {
		return err
	}
	// Convert responses
	hc.Id = shadow.Id
	hc.Type = shadow.Type
	hc.CreatedBy = shadow.CreatedBy
	hc.UpdatedBy = shadow.UpdatedBy
	hc.Number = shadow.Number
	hc.Title = shadow.Title
	hc.Description = shadow.Description
	hc.Tags = shadow.Tags
	hc.Flag = shadow.Flag
	hc.Tlp = shadow.Tlp
	hc.TlpLabel = shadow.TlpLabel
	hc.Pap = shadow.Pap
	hc.PapLabel = shadow.PapLabel
	hc.Status = shadow.Status
	hc.Stage = shadow.Stage
	hc.Summary = shadow.Summary
	hc.Severity = shadow.Severity
	hc.ImpactStatus = shadow.ImpactStatus
	hc.Assignee = shadow.Assignee
	hc.CustomFields = shadow.CustomFields
	hc.UserPermissions = shadow.UserPermissions
	hc.ExtraData = shadow.ExtraData

	// For int64 fields, use the convertInt64ToTime function:
	hc.CreatedAt = convertInt64ToTime(shadow.CreatedAt)
	hc.UpdatedAt = convertInt64ToTime(shadow.UpdatedAt)
	hc.StartDate = convertInt64ToTime(shadow.StartDate)
	hc.EndDate = convertInt64ToTime(shadow.EndDate)
	hc.NewDate = convertInt64ToTime(shadow.NewDate)
	hc.InProgressDate = convertInt64ToTime(shadow.InProgressDate)
	hc.ClosedDate = convertInt64ToTime(shadow.ClosedDate)
	hc.AlertDate = convertInt64ToTime(shadow.AlertDate)
	hc.AlertNewDate = convertInt64ToTime(shadow.AlertNewDate)
	hc.AlertInProgressDate = convertInt64ToTime(shadow.AlertInProgressDate)
	hc.AlertImportedDate = convertInt64ToTime(shadow.AlertImportedDate)
	hc.TimeToDetect = convertInt64ToDuration(shadow.TimeToDetect)
	hc.TimeToTriage = convertInt64ToDuration(shadow.TimeToTriage)
	hc.TimeToQualify = convertInt64ToDuration(shadow.TimeToQualify)
	hc.TimeToAcknowledge = convertInt64ToDuration(shadow.TimeToAcknowledge)
	hc.TimeToResolve = convertInt64ToDuration(shadow.TimeToResolve)
	hc.HandlingDuration = convertInt64ToDuration(shadow.HandlingDuration)
	return nil

}

// HiveSearch is used to create the root search query
type HiveSearch struct {
	Query []SearchQuery `json:"query"`
}

// SearchQuery includes all available filter for the query api endpoint on thehive5
type SearchQuery struct {
	Name       string                `json:"_name"`
	And        *[]Filter             `json:"_and,omitempty"`
	Or         *[]Filter             `json:"_or,omitempty"`
	Any        *Filter               `json:"_any,omitempty"` // not properly documented
	Not        *Filter               `json:"_not,omitempty"` // not properly documented
	Lt         *Filter               `json:"_lt,omitempty"`
	Gt         *Filter               `json:"_gt,omitempty"`
	Lte        *Filter               `json:"_lte,omitempty"`
	Gte        *Filter               `json:"_gte,omitempty"`
	Ne         *Filter               `json:"_ne,omitempty"`
	Eq         *Filter               `json:"_eq,omitempty"`
	Is         *Filter               `json:"_is,omitempty"`
	StartsWith *Filter               `json:"_startsWith,omitempty"`
	EndsWith   *Filter               `json:"_endsWith,omitempty"`
	Id         *string               `json:"_id,omitempty"`
	Between    *Filter               `json:"_between,omitempty"`
	In         *Filter               `json:"_in,omitempty"`
	Contains   *string               `json:"_contains,omitempty"`
	Like       *Filter               `json:"_like,omitempty"`
	Match      *Filter               `json:"_match,omitempty"`
	Sort       *[1]map[string]string `json:"_fields,omitempty"`
	ScopeFrom  *int                  `json:"from,omitempty"`
	ScopeTo    *int                  `json:"to,omitempty"`
	IdOrName   *string               `json:"idOrName,omitempty"`
}

// A CaseStatusResponse is used for containing all possible case status options
type CaseStatusResponse struct {
	Id          string            `json:"_id"`
	Type        string            `json:"_type"`
	UpdatedAt   time.Time         `json:"_updatedAt,omitempty"`
	UpdatedBy   string            `json:"_updatedBy,omitempty"`
	CreatedAt   time.Time         `json:"_createdAt"`
	CreatedBy   string            `json:"_createdBy"`
	Value       string            `json:"value"`
	Stage       string            `json:"stage"`
	Order       int               `json:"order,omitempty"`
	Description string            `json:"description,omitempty"`
	Colour      string            `json:"colour,omitempty"`
	ExtraData   map[string]string `json:"extraData"`
}

// A CaseStatusResponse is used for containing all possible case status options
type shadowCaseStatusResponse struct {
	Id          string            `json:"_id"`
	Type        string            `json:"_type"`
	UpdatedAt   int64             `json:"_updatedAt"`
	UpdatedBy   string            `json:"_updatedBy"`
	CreatedAt   int64             `json:"_createdAt"`
	CreatedBy   string            `json:"_createdBy"`
	Value       string            `json:"value"`
	Stage       string            `json:"stage"`
	Order       int               `json:"order"`
	Description string            `json:"description"`
	Colour      string            `json:"colour"`
	ExtraData   map[string]string `json:"extraData"`
}

func (cs *CaseStatusResponse) UnmarshalJSON(data []byte) error {
	shadow := new(shadowCaseStatusResponse)
	err := json.Unmarshal(data, &shadow)
	if err != nil {
		return err
	}

	cs.Id = shadow.Id
	cs.Type = shadow.Type
	cs.CreatedAt = convertInt64ToTime(shadow.CreatedAt)
	cs.Value = shadow.Value
	cs.Stage = shadow.Stage
	cs.Order = shadow.Order
	cs.Colour = shadow.Colour
	cs.ExtraData = shadow.ExtraData

	return nil
}

// A Filter is used for filtering on the query endpoint
type Filter struct {
	Field     string            `json:"_field"`
	Fields    map[string]string `json:"-,omitempty"`
	Value     interface{}       `json:"_value,omitempty"`
	Values    []interface{}     `json:"_values,omitempty"`
	Scope     *Scope            `json:",omitempty"`
	ExtraData *struct{}         `json:"extraData,omitempty"`
}

// A Scope has to be defined to search on specific pages
type Scope struct {
	From *int `json:"from,omitempty"`
	To   *int `json:"to,omitempty"`
}

// executeCaseSearchQuery is a helper function to do query related searches
func (hive *Hivedata) executeCaseSearchQuery(query []byte) ([]HiveCaseResponse, error) {
	url, err := url.JoinPath(hive.Url, "/api/v1/query")
	if err != nil {
		return nil, err
	}

	ret, err := hive.webRequest(url, POST, query)
	if err != nil {
		return nil, err
	}

	var parsedRet []HiveCaseResponse
	err = json.Unmarshal(ret, &parsedRet)
	return parsedRet, err
}

// executeCaseStatusQuery is a helper function to do query related searches
func (hive *Hivedata) executeCaseStatusQuery(query []byte) ([]CaseStatusResponse, error) {
	url, err := url.JoinPath(hive.Url, "/api/v1/query")
	if err != nil {
		return nil, err
	}

	ret, err := hive.webRequest(url, POST, query)
	if err != nil {
		return nil, err
	}

	var parsedRet []CaseStatusResponse
	err = json.Unmarshal(ret, &parsedRet)
	return parsedRet, err
}

// GetCaseStatusOptions returns all options that are able to be set on a case
func (hive *Hivedata) GetCaseStatusOptions() ([]CaseStatusResponse, error) {
	query, err := hive.createSearchQuery(
		SearchQuery{Name: "listCaseStatus"},
		SearchQuery{Name: "sort", Sort: &[1]map[string]string{{"stage": "desc"}}},
	)
	if err != nil {
		return nil, err
	}

	return hive.executeCaseStatusQuery(query)
}

// FindCaseByCustomField finds cases acoording to the query values on a specific custom field submitted
func (hive *Hivedata) FindCaseByCustomField(queryfield string, queryvalue string) ([]HiveCaseResponse, error) {
	query, err := hive.createSearchQuery(
		SearchQuery{Name: "listCase"},
		SearchQuery{Name: "filter", Eq: &Filter{Field: fmt.Sprintf("customFields.%s", strings.ToLower(queryfield)), Value: &queryvalue}},
		SearchQuery{Name: "sort", Sort: &[1]map[string]string{{"_updatedAt": "desc"}}},
	)
	if err != nil {
		return nil, err
	}

	return hive.executeCaseSearchQuery(query)
}

// FindCase allows to search for self defined case queries
func (hive *Hivedata) FindCase(searchQuery []SearchQuery) ([]HiveCaseResponse, error) {

	query, err := hive.createSearchQuery(searchQuery...)
	if err != nil {
		return nil, err
	}

	return hive.executeCaseSearchQuery(query)
}

func (hive *Hivedata) DeleteCase(caseId int) error {
	caseNumber := strconv.Itoa(caseId)
	url, err := url.JoinPath(hive.Url, "/api/v1/case", caseNumber)
	if err != nil {
		return err
	}

	_, err = hive.webRequest(url, DELETE, nil)
	return err
}

// CreateCase is used to add a new case on thehive5
// Returns HiveCase struct and response error
func (hive *Hivedata) CreateCase(newCase *HiveCase) (*HiveCaseResponse, error) {
	url, err := url.JoinPath(hive.Url, "/api/v1/case")
	if err != nil {
		return nil, err
	}
	jsondata, err := json.Marshal(newCase)
	if err != nil {
		return nil, err
	}

	ret, err := hive.webRequest(url, POST, jsondata)
	if err != nil {
		return nil, err
	}

	var parsedRet HiveCaseResponse
	err = json.Unmarshal(ret, &parsedRet)
	if err != nil {
		return nil, err
	}

	return &parsedRet, err
}

// UpdateCase is used to update an existing case on thehive5
// Only returns data if an error occurs
func (hive *Hivedata) UpdateCase(idOrName int, updatedCase *HiveUpdateCase) error {
	url, err := url.JoinPath(hive.Url, "/api/v1/case", strconv.Itoa(idOrName))
	if err != nil {
		return err
	}

	jsondata, err := json.Marshal(updatedCase)
	if err != nil {
		return err
	}

	_, err = hive.webRequest(url, PATCH, jsondata)
	return err
}

// CreateCaseFromAlert creates a new case from an existing alert
// Returns newly created case
func (hive *Hivedata) CreateCaseFromAlert(alertId string, alert *HiveCase) (*HiveCaseResponse, error) {
	url, err := url.JoinPath(hive.Url, "/api/v1/alert/", alertId, "/case")
	if err != nil {
		return nil, err
	}

	jsondata, err := json.Marshal(alert)
	if err != nil {
		return nil, err
	}

	ret, err := hive.webRequest(url, POST, jsondata)
	if err != nil {
		return nil, err
	}

	parsedRet := HiveCaseResponse{}
	err = json.Unmarshal(ret, &parsedRet)
	if err != nil {
		return nil, err
	}
	return &parsedRet, err
}

// GetCase looks up a case by ID and returns it
func (hive *Hivedata) GetCase(caseId int) (*HiveCaseResponse, error) {
	caseNumber := strconv.Itoa(caseId)
	url, err := url.JoinPath(hive.Url, "/api/v1/case/", caseNumber)
	if err != nil {
		return nil, err
	}

	ret, err := hive.webRequest(url, GET, nil)
	if err != nil {
		return nil, err
	}

	var parsedRet *HiveCaseResponse
	err = json.Unmarshal(ret, &parsedRet)
	return parsedRet, err
}

// GetCasesTimed takes time object to look back for a certain time and returns all the cases found
// timeframe is always timeframe < xxx which means that everything since the timeframe will be returned
func (hive *Hivedata) GetCasesTimed(timeframe time.Time) ([]HiveCaseResponse, error) {
	// hive expects a miliseconds time string
	time := strconv.FormatInt(timeframe.UnixMilli(), 10)

	query, err := hive.createSearchQuery(
		SearchQuery{Name: "listCase"},
		SearchQuery{Name: "filter", Gte: &Filter{Field: "_updatedAt", Value: &time}})
	if err != nil {
		return nil, err
	}

	return hive.executeCaseSearchQuery(query)
}

// GetCaseAlerts returns all alerts associated with a case
// It returns a AlertsResponse slice or an error
func (hive *Hivedata) GetCaseAlerts(caseId int) ([]HiveAlertResponse, error) {
	caseNumber := strconv.Itoa(caseId)
	query, err := hive.createSearchQuery(
		SearchQuery{Name: "getCase", IdOrName: &caseNumber},
		SearchQuery{Name: "alerts"})
	if err != nil {
		return nil, err
	}

	return hive.executeAlertSearchQuery(query)
}
