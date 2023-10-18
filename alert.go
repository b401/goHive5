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

// A HiveAlert stores an alert for updating and creating new alerts on thehive5
type HiveAlert struct {
	Type         string         `json:"type"`
	Source       string         `json:"source"`
	SourceRef    string         `json:"sourceRef"`
	Title        string         `json:"title"`
	Description  string         `json:"description"`
	Severity     string         `json:"severity,omitempty"`
	Date         time.Time      `json:"date,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	ExternalLink string         `json:"externalLink,omitempty"`
	Flag         bool           `json:"flag,omitempty"`
	Tlp          string         `json:"tlp,omitempty"`
	Pap          string         `json:"pap,omitempty"`
	CustomFields *[]CustomField `json:"customFields,omitempty"`
	Summary      string         `json:"summary,omitempty"`
	Status       string         `json:"status,omitempty"`
	Assignee     string         `json:"assignee,omitempty"`
	CaseTemplate string         `json:"caseTemplate,omitempty"`
	Observables  *[]Observable  `json:"observables,omitempty"`
	Procedures   *[]Procedure   `json:"procedures,omitempty"`
}

// A HiveUpdateAlert stores information to update an existing alert
// Looks incomplete on theHive api documentation
type HiveUpdateAlert struct {
	Type         string        `json:"type,omitempty"`
	Source       string        `json:"source,omitempty"`
	SourceRef    string        `json:"sourceRef,omitempty"`
	ExternalLink string        `json:"externalLink,omitempty"`
	Title        string        `json:"title,omitempty"`
	Description  string        `json:"description,omitempty"`
	Severity     string        `json:"severity,omitempty"`
	Date         time.Time     `json:"date,omitempty"`
	LastSyncDate time.Time     `json:"lastSyncDate,omitempty"`
	Tags         []string      `json:"tags,omitempty"`
	Tlp          string        `json:"tlp,omitempty"`
	Pap          string        `json:"pap,omitempty"`
	Follow       *bool         `json:"follow,omitempty"`
	CustomFields []CustomField `json:"customFields,omitempty"`
	Status       string        `json:"status,omitempty"`
	Summary      string        `json:"summary,omitempty"`
	Assignee     string        `json:"assignee,omitempty"`
	AddTags      []string      `json:"addTags,omitempty"`
	RemoveTags   []string      `json:"removeTags,omitempty"`
}

// An HiveAlertResponse contains the attributes thehive5 sends back on requests
// fields may contain nil
type HiveAlertResponse struct {
	Id                string        `json:"_id"`
	Type              string        `json:"_type"`
	CreatedBy         string        `json:"_createdBy"`
	UpdatedBy         string        `json:"_updatedBy"`
	CreatedAt         time.Time     `json:"_createdAt"`
	UpdatedAt         time.Time     `json:"_updatedAt"`
	AlertType         string        `json:"type"`
	Source            string        `json:"source"`
	SourceRef         string        `json:"sourceRef"`
	ExternalLink      string        `json:"externalLink"`
	Title             string        `json:"title"`
	Description       string        `json:"description"`
	Severity          int           `json:"severity"`
	SeverityLabel     string        `json:"severityLabel"`
	Date              time.Time     `json:"date"`
	Tags              []string      `json:"tags"`
	Tlp               int           `json:"tlp"`
	TlpLabel          string        `json:"tlpLabel"`
	Pap               int           `json:"pap"`
	PapLabel          string        `json:"papLabel"`
	Follow            bool          `json:"follow"`
	CustomFields      []CustomField `json:"customFields"`
	CaseTemplate      string        `json:"caseTemplate"`
	ObservableCount   int64         `json:"observableCount"`
	CaseID            string        `json:"caseId"`
	Status            string        `json:"status"`
	Stage             string        `json:"stage"`
	Assignee          string        `json:"assignee"`
	Summary           string        `json:"summary"`
	ExtraData         struct{}      `json:"extraData"`
	NewDate           time.Time     `json:"newDate"`
	InProgressDate    time.Time     `json:"inProgressDate"`
	ClosedDate        time.Time     `json:"closedDate"`
	ImportedDate      time.Time     `json:"importedDate"`
	TimeToDetect      time.Duration `json:"timeToDetect"`
	TimeToTriage      time.Duration `json:"timeToTriage"`
	TimeToQualify     time.Duration `json:"timeToQualify"`
	TimeToAcknowledge time.Duration `json:"timeToAcknowledge"`
}

// shadowHiveAlertResponse is used to unmarshall the int64 values int time.Time and time.Duration
type shadowHiveAlertResponse struct {
	Id                string        `json:"_id"`
	Type              string        `json:"_type"`
	CreatedBy         string        `json:"_createdBy"`
	UpdatedBy         string        `json:"_updatedBy"`
	CreatedAt         int64         `json:"_createdAt"`
	UpdatedAt         int64         `json:"_updatedAt"`
	AlertType         string        `json:"type"`
	Source            string        `json:"source"`
	SourceRef         string        `json:"sourceRef"`
	ExternalLink      string        `json:"externalLink"`
	Title             string        `json:"title"`
	Description       string        `json:"description"`
	Severity          int           `json:"severity"`
	SeverityLabel     string        `json:"severityLabel"`
	Date              int64         `json:"date"`
	Tags              []string      `json:"tags"`
	Tlp               int           `json:"tlp"`
	TlpLabel          string        `json:"tlpLabel"`
	Pap               int           `json:"pap"`
	PapLabel          string        `json:"papLabel"`
	Follow            bool          `json:"follow"`
	CustomFields      []CustomField `json:"customFields"`
	CaseTemplate      string        `json:"caseTemplate"`
	ObservableCount   int64         `json:"observableCount"`
	CaseID            string        `json:"caseId"`
	Status            string        `json:"status"`
	Stage             string        `json:"stage"`
	Assignee          string        `json:"assignee"`
	Summary           string        `json:"summary"`
	ExtraData         struct{}      `json:"extraData"`
	NewDate           int64         `json:"newDate"`
	InProgressDate    int64         `json:"inProgressDate"`
	ClosedDate        int64         `json:"closedDate"`
	ImportedDate      int64         `json:"importedDate"`
	TimeToDetect      int64         `json:"timeToDetect"`
	TimeToTriage      int64         `json:"timeToTriage"`
	TimeToQualify     int64         `json:"timeToQualify"`
	TimeToAcknowledge int64         `json:"timeToAcknowledge"`
}

// Unmarshal thehive5 returned values into the HiveAlertResponse structs. Making sure that int64 gets converted into time.Time
func (ar *HiveAlertResponse) UnmarshalJSON(data []byte) error {
	shadow := new(shadowHiveAlertResponse)
	err := json.Unmarshal(data, &shadow)
	if err != nil {
		return err
	}

	// Convert responses
	ar.Id = shadow.Id
	ar.Type = shadow.Type
	ar.CreatedBy = shadow.CreatedBy
	ar.UpdatedBy = shadow.UpdatedBy
	ar.CreatedAt = convertInt64ToTime(shadow.CreatedAt)
	ar.UpdatedAt = convertInt64ToTime(shadow.UpdatedAt)
	ar.AlertType = shadow.AlertType
	ar.Source = shadow.Source
	ar.SourceRef = shadow.SourceRef
	ar.ExternalLink = shadow.ExternalLink
	ar.Title = shadow.Title
	ar.Description = shadow.Description
	ar.Severity = shadow.Severity
	ar.SeverityLabel = shadow.SeverityLabel
	ar.Date = convertInt64ToTime(shadow.Date)
	ar.Tags = shadow.Tags
	ar.Tlp = shadow.Tlp
	ar.TlpLabel = shadow.TlpLabel
	ar.Pap = shadow.Pap
	ar.PapLabel = shadow.PapLabel
	ar.Follow = shadow.Follow
	ar.CustomFields = shadow.CustomFields
	ar.CaseTemplate = shadow.CaseTemplate
	ar.ObservableCount = shadow.ObservableCount
	ar.CaseID = shadow.CaseID
	ar.Status = shadow.Status
	ar.Stage = shadow.Stage
	ar.Assignee = shadow.Assignee
	ar.Summary = shadow.Summary
	ar.ExtraData = shadow.ExtraData
	ar.NewDate = convertInt64ToTime(shadow.NewDate)
	ar.InProgressDate = convertInt64ToTime(shadow.InProgressDate)
	ar.ClosedDate = convertInt64ToTime(shadow.ClosedDate)
	ar.ImportedDate = convertInt64ToTime(shadow.ImportedDate)
	ar.TimeToDetect = convertInt64ToDuration(shadow.TimeToDetect)
	ar.TimeToTriage = convertInt64ToDuration(shadow.TimeToTriage)
	ar.TimeToQualify = convertInt64ToDuration(shadow.TimeToQualify)
	ar.TimeToAcknowledge = convertInt64ToDuration(shadow.TimeToAcknowledge)

	return nil

}

// Marshalling the alert requests
func (ha *HiveAlert) MarshalJSON() ([]byte, error) {
	type Alias HiveAlert
	var (
		dateInt64   int64
		tlpInt      int
		papInt      int
		severityInt int
	)

	// We set the type if the caller hasn't done it
	if len(ha.Type) == 0 {
		ha.Type = "alert"
	}

	if len(ha.Tlp) != 0 {
		var tlp Tlp
		err := tlp.FromString(ha.Tlp)
		if err != nil {
			return nil, err
		}
		tlpInt = int(tlp)
	}
	if len(ha.Pap) != 0 {
		var pap Pap
		err := pap.FromString(ha.Pap)
		if err != nil {
			return nil, err
		}
		papInt = int(pap)
	}

	if len(ha.Severity) != 0 {
		var sev Severity
		err := sev.FromString(ha.Severity)
		if err != nil {
			return nil, err
		}
		severityInt = int(sev)
	}

	if !ha.Date.IsZero() {
		dateInt64 = ha.Date.UTC().UnixMilli()
	}

	return json.Marshal(&struct {
		Date     int64  `json:"date,omitempty"`
		Type     string `json:"type"`
		Tlp      int    `json:"tlp,omitempty"`
		Pap      int    `json:"pap,omitempty"`
		Severity int    `json:"severity,omitempty"`
		*Alias
	}{
		Type:     ha.Type,
		Date:     dateInt64,
		Tlp:      tlpInt,
		Pap:      papInt,
		Severity: severityInt,
		Alias:    (*Alias)(ha),
	})

}

// Marshalling the alert update requests
func (ha *HiveUpdateAlert) MarshalJSON() ([]byte, error) {
	type Alias HiveUpdateAlert
	var (
		dateInt64    int64
		lastSyncDate int64
		tlpInt       int
		papInt       int
		severityInt  int
	)

	// We set the type if the caller hasn't done it
	if len(ha.Type) == 0 {
		ha.Type = "alert"
	}

	if len(ha.Tlp) != 0 {
		var tlp Tlp
		err := tlp.FromString(ha.Tlp)
		if err != nil {
			return nil, err
		}
		tlpInt = int(tlp)
	}
	if len(ha.Pap) != 0 {
		var pap Pap
		err := pap.FromString(ha.Pap)
		if err != nil {
			return nil, err
		}
		papInt = int(pap)
	}

	if len(ha.Severity) != 0 {
		var sev Severity
		err := sev.FromString(ha.Severity)
		if err != nil {
			return nil, err
		}
		severityInt = int(sev)
	}

	if !ha.Date.IsZero() {
		dateInt64 = ha.Date.UTC().UnixMilli()
	}

	if !ha.LastSyncDate.IsZero() {
		lastSyncDate = ha.LastSyncDate.UTC().UnixMilli()
	}

	return json.Marshal(&struct {
		Date         int64  `json:"date,omitempty"`
		LastSyncDate int64  `json:"lastSyncDate,omitempty"`
		Type         string `json:"type"`
		Tlp          int    `json:"tlp,omitempty"`
		Pap          int    `json:"pap,omitempty"`
		Severity     int    `json:"severity,omitempty"`
		*Alias
	}{
		Type:         ha.Type,
		Date:         dateInt64,
		LastSyncDate: lastSyncDate,
		Tlp:          tlpInt,
		Pap:          papInt,
		Severity:     severityInt,
		Alias:        (*Alias)(ha),
	})
}

// A CustomField contains the custom field declaration on thehive5
// Make sure that the name attribute exists on your thehive5 instance
type CustomField struct {
	Name        string             `json:"name"`
	DisplayName string             `json:"displayName,omitempty"`
	Group       string             `json:"group"`
	Value       interface{}        `json:"value"`
	Description string             `json:"description"`
	Type        string             `json:"type"`
	Mandatory   bool               `json:"mandatory,omitempty"`
	Options     *map[string]string `json:"options,omitempty"`
}

// MarshalJSON sets the customField name to lowercase as all customFields on thehive5 are in that format as well
func (c *CustomField) MarshalJSON() ([]byte, error) {
	type Alias CustomField

	// we try to automatically detect the type
	if len(c.Type) == 0 {
		detectedFieldType, err := detectCustomFieldType(c.Value)
		if err != nil {
			return nil, err
		}
		c.Type = *detectedFieldType
	}

	// we try to convert for you
	if c.Type == "string" {
		c.Value = fmt.Sprintf("%s", c.Value)
	}

	return json.Marshal(&struct {
		Name  string `json:"name"`
		Value string `json:"value"`
		*Alias
	}{
		Name:  strings.ToLower(c.Name),
		Value: c.Value.(string),
		Alias: (*Alias)(c),
	})
}

type CustomFieldResponse struct {
	Id          string              `json:"_id"`
	Type        string              `json:"_type,omitempty"`
	CreatedBy   string              `json:"_createdBy,omitempty"`
	UpdatedBy   string              `json:"_updatedBy"`
	CreatedAt   time.Time           `json:"_createdAt"`
	UpdatedAt   time.Time           `json:"_updatedAt,omitempty"`
	Name        string              `json:"name"`
	DisplayName string              `json:"displayName"`
	Group       string              `json:"group"`
	Description string              `json:"description"`
	FieldType   string              `json:"type"`
	Options     []map[string]string `json:"options,omitempty"`
	Mandatory   bool                `json:"mandatory"`
	ExtraData   map[string]string   `json:"extraData"`
}

type shadowCustomFieldResponse struct {
	Id          string              `json:"_id"`
	Type        string              `json:"_type,omitempty"`
	CreatedBy   string              `json:"_createdBy,omitempty"`
	UpdatedBy   string              `json:"_updatedBy"`
	CreatedAt   int64               `json:"_createdAt"`
	UpdatedAt   int64               `json:"_updatedAt,omitempty"`
	Name        string              `json:"name"`
	DisplayName string              `json:"displayName"`
	Group       string              `json:"group"`
	Description string              `json:"description"`
	FieldType   string              `json:"type"`
	Options     []map[string]string `json:"options,omitempty"`
	Mandatory   bool                `json:"mandatory"`
	ExtraData   map[string]string   `json:"extraData"`
}

// Unmarshal thehive5 returned values into the HiveAlertResponse structs. Making sure that int64 gets converted into time.Time
func (c *CustomFieldResponse) UnmarshalJSON(data []byte) error {
	shadow := new(shadowCustomFieldResponse)
	err := json.Unmarshal(data, &shadow)
	if err != nil {
		return err
	}

	c.Id = shadow.Id
	c.Type = shadow.Type
	c.CreatedBy = shadow.CreatedBy
	c.UpdatedBy = shadow.UpdatedBy
	c.CreatedAt = convertInt64ToTime(shadow.CreatedAt)
	c.UpdatedAt = convertInt64ToTime(shadow.UpdatedAt)
	c.Name = shadow.Name
	c.DisplayName = shadow.DisplayName
	c.Group = shadow.Group
	c.Description = shadow.Description
	c.FieldType = shadow.FieldType
	c.Options = shadow.Options
	c.Mandatory = shadow.Mandatory
	c.ExtraData = shadow.ExtraData

	return nil
}

// executeAlertSearchQuery is a helper function to do query related searches
func (hive *Hivedata) executeAlertSearchQuery(query []byte) ([]HiveAlertResponse, error) {
	url, err := url.JoinPath(hive.Url, "/api/v1/query")
	if err != nil {
		return nil, err
	}
	ret, err := hive.webRequest(url, POST, query)

	if err != nil {
		return nil, err
	}

	var parsedRet []HiveAlertResponse
	err = json.Unmarshal(ret, &parsedRet)
	return parsedRet, err
}

// FindAlertsByFieldTimed allows a lookback for a specific time for the _UpdatedArt field and a specific field & value
func (hive *Hivedata) FindAlertsByFieldTimed(queryfield string, queryvalue string, timeframe time.Time) ([]HiveAlertResponse, error) {
	// hive expects a milliseconds time string
	time := strconv.FormatInt(timeframe.Unix()*1000, 10)

	query, err := hive.createSearchQuery(
		SearchQuery{Name: "listAlert"},
		SearchQuery{Name: "filter", Eq: &Filter{Field: strings.ToLower(queryfield), Value: &queryvalue}},
		SearchQuery{Name: "filter", Gte: &Filter{Field: "_updatedAt", Value: &time}},
	)
	if err != nil {
		return nil, err
	}

	return hive.executeAlertSearchQuery(query)
}

// GetAlertsTimed returns all alerts which were updated since a specific date
func (hive *Hivedata) GetAlertsTimed(timeframe time.Time) ([]HiveAlertResponse, error) {
	time := strconv.FormatInt(timeframe.Unix()*1000, 10)

	query, err := hive.createSearchQuery(
		SearchQuery{Name: "listAlert"},
		SearchQuery{Name: "filter", Gte: &Filter{Field: "_updatedAt", Value: &time}},
	)
	if err != nil {
		return nil, err
	}

	return hive.executeAlertSearchQuery(query)
}

// FindAlertsByCustomField does the same thing as FindAlertsByField but for custom fields.
// Use this function for custom fields
func (hive *Hivedata) FindAlertsByCustomField(queryfield string, queryvalue string) ([]HiveAlertResponse, error) {

	// Creates the json struct object
	query, err := hive.createSearchQuery(
		SearchQuery{Name: "listAlert"},
		SearchQuery{Name: "filter", Eq: &Filter{Field: fmt.Sprintf("customFields.%s", strings.ToLower(queryfield)), Value: &queryvalue}},
	)

	if err != nil {
		return nil, err
	}

	return hive.executeAlertSearchQuery(query)
}

// MergeAlert merges an alert into a case.
// The alertId must be a string, while the caseNumber must be an int.
// It returns an error if the merging process fails.
func (hive *Hivedata) MergeAlert(alertId string, caseNumber int) error {
	url, err := url.JoinPath(hive.Url, "/api/v1/alert/", alertId, "/merge/", strconv.Itoa(caseNumber))
	if err != nil {
		return err
	}

	_, err = hive.webRequest(url, POST, nil)
	return err
}

// CreateAlert adds a new alert on thehive5 and returns the created alert response.
func (hive *Hivedata) CreateAlert(alertObject *HiveAlert) (*HiveAlertResponse, error) {
	url, err := url.JoinPath(hive.Url, "api/v1/alert")
	if err != nil {
		return nil, err
	}
	jsondata, err := json.Marshal(alertObject)
	if err != nil {
		return nil, err
	}

	ret, err := hive.webRequest(url, POST, jsondata)
	if err != nil {
		return nil, err
	}

	parsedRet := new(HiveAlertResponse)
	err = json.Unmarshal(ret, parsedRet)
	if err != nil {
		return nil, err
	}

	return parsedRet, err
}

// DeleteAlert deletes an alert and returns an error if the deletion fails.
func (hive *Hivedata) DeleteAlert(alertId string) error {
	url, err := url.JoinPath(hive.Url, "api/v1/alert", alertId)
	if err != nil {
		return err
	}
	_, err = hive.webRequest(url, DELETE, nil)
	return err
}

// GetAlert retrieves a single Alert using its alertId.
// It returns an HiveAlertResponse or an error.
func (hive *Hivedata) GetAlert(alertId string) (*HiveAlertResponse, error) {
	url, err := url.JoinPath(hive.Url, "api/v1/alert", alertId)
	if err != nil {
		return nil, err
	}

	ret, err := hive.webRequest(url, GET, nil)
	if err != nil {
		return nil, err
	}

	parsedRet := new(HiveAlertResponse)
	err = json.Unmarshal(ret, parsedRet)
	return parsedRet, err
}

// UpdateAlert updates an alert given a HiveUpdateAlert struct
func (hive *Hivedata) UpdateAlert(alertId string, alert *HiveUpdateAlert) error {
	url, err := url.JoinPath(hive.Url, "api/v1/alert", alertId)
	if err != nil {
		return err
	}

	jsondata, err := json.Marshal(alert)
	if err != nil {
		return err
	}
	_, err = hive.webRequest(url, PATCH, jsondata)
	return err
}

// AddAlertObservable adds a new observable to an existing alert.
// Returns an error if the addition fails.
func (hive *Hivedata) AddAlertObservable(alertNumber string, observable Observable) error {
	url, err := url.JoinPath(hive.Url, "api/v1/alert", alertNumber, "observable")
	if err != nil {
		return err
	}

	jsondata, err := json.Marshal(observable)
	if err != nil {
		return err
	}

	_, err = hive.webRequest(url, POST, jsondata)
	return err
}

// GetAlertObservables returns all observables associated with an alert
func (hive *Hivedata) GetAlertObservables(alertId string) ([]ObservableResponse, error) {
	query, err := hive.createSearchQuery(
		SearchQuery{Name: "getAlert", IdOrName: alertId},
		SearchQuery{Name: "observables"},
	)

	if err != nil {
		return nil, err
	}

	return hive.executeObservableSearchQuery(query)
}

// GetAlertObservable returns a single observable associated with an alert.
// Use this if you need to get all alerts that have a specific observable with a specific value
// Example: queryfield: data, queryvalue: 127.0.0.1
func (hive *Hivedata) GetAlertObservable(alertId, queryfield, queryvalue string) (*ObservableResponse, error) {

	query, err := hive.createSearchQuery(
		SearchQuery{Name: "getAlert", IdOrName: alertId},
		SearchQuery{Name: "observables"},
		SearchQuery{Name: "filter", Eq: &Filter{Field: queryfield, Value: &queryvalue}},
	)

	if err != nil {
		return nil, err
	}

	resp, err := hive.executeObservableSearchQuery(query)
	if err != nil {
		return nil, err
	}

	if len(resp) > 0 {
		return &resp[0], err
	}

	return nil, fmt.Errorf("no observable found")
}
