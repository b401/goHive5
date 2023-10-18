/*
thehive5 implements functionality to interact with the most recent version of thehive.
https://www.strangebee.com/thehive/
*/
package thehive5

import (
	"encoding/json"
	"net/url"
	"os"
	"strconv"
	"time"
)

// An Observable is used to define objects that have been seen in an alert on an incident
// It was renamed from Artifact to Observable in thehive5
type Observable struct {
	DataType         string    `json:"dataType,omitempty"`
	Data             string    `json:"data,omitempty"`
	Message          string    `json:"message,omitempty"`
	Tlp              string    `json:"tlp,omitempty"`
	Pap              string    `json:"pap,omitempty"`
	Tags             []string  `json:"tags,omitempty"`
	Ioc              bool      `json:"ioc,omitempty"`
	Sighted          bool      `json:"sighted,omitempty"`
	StartDate        time.Time `json:"startDate,omitempty"`
	SightedAt        time.Time `json:"sightedAt,omitempty"`
	IgnoreSimilarity bool      `json:"ignoreSimilarity,omitempty"`
	IsZip            bool      `json:"isZip,omitempty"`
	ZipPassword      string    `json:"zipPassword,omitempty"`
}

// Marshalling the observables
func (o *Observable) MarshalJSON() ([]byte, error) {
	type Alias Observable
	var (
		startDateInt64 *int64
		sightedAtInt64 *int64
		tlpInt         *int
		papInt         *int
	)

	if len(o.Tlp) != 0 {
		var tlp Tlp
		err := tlp.FromString(o.Tlp)
		if err != nil {
			return nil, err
		}
		tmp := int(tlp)
		tlpInt = &tmp
	}
	if len(o.Pap) != 0 {
		var pap Pap
		err := pap.FromString(o.Pap)
		if err != nil {
			return nil, err
		}
		tmp := int(pap)
		papInt = &tmp
	}

	if !o.StartDate.UTC().IsZero() {
		dateTmp := o.StartDate.UTC().UnixMilli()
		startDateInt64 = &dateTmp
	}
	if !o.SightedAt.UTC().IsZero() {
		dateTmp := o.SightedAt.UTC().UnixMilli()
		sightedAtInt64 = &dateTmp
		sighted := true
		o.Sighted = sighted
	}

	return json.Marshal(&struct {
		StartDate *int64 `json:"startDate,omitempty"`
		SightedAt *int64 `json:"sightedAt,omitempty"`
		Tlp       *int   `json:"tlp,omitempty"`
		Pap       *int   `json:"pap,omitempty"`
		*Alias
	}{
		StartDate: startDateInt64,
		SightedAt: sightedAtInt64,
		Tlp:       tlpInt,
		Pap:       papInt,
		Alias:     (*Alias)(o),
	})
}

// ObservableResponse contains the returned values from thehive5 observable api
type ObservableResponse struct {
	Id               string     `json:"_id"`
	Type             string     `json:"_type"`
	CreatedBy        string     `json:"_createdBy"`
	UpdatedBy        string     `json:"_updatedBy"`
	CreatedAt        time.Time  `json:"_createdAt"`
	UpdatedAt        time.Time  `json:"_updatedAt"`
	DataType         string     `json:"dataType"`
	Data             string     `json:"data"`
	StartDate        time.Time  `json:"startDate"`
	Attachment       Attachment `json:"attachment"`
	Tlp              int        `json:"tlp"`
	TlpLabel         string     `json:"tlpLabel"`
	Pap              int        `json:"pap"`
	PapLabel         string     `json:"papLabel"`
	Tags             []string   `json:"tags"`
	Ioc              bool       `json:"ioc"`
	Sighted          bool       `json:"sighted"`
	SightedAt        time.Time  `json:"sightedAt"`
	Reports          struct{}   `json:"reports"`
	Message          string     `json:"message"`
	ExtraData        ExtraData `json:"extraData,omitempty"`
	IgnoreSimilarity bool       `json:"ignoreSimilarity"`
}

// ObservableResponse contains the returned values from thehive5 observable api
type shadowObservableResponse struct {
	Id               string     `json:"_id"`
	Type             string     `json:"_type"`
	CreatedBy        string     `json:"_createdBy"`
	UpdatedBy        string     `json:"_updatedBy"`
	CreatedAt        int64      `json:"_createdAt"`
	UpdatedAt        int64      `json:"_updatedAt"`
	DataType         string     `json:"dataType"`
	Data             string     `json:"data"`
	StartDate        int64      `json:"startDate"`
	Attachment       Attachment `json:"attachment"`
	Tlp              int        `json:"tlp"`
	TlpLabel         string     `json:"tlpLabel"`
	Pap              int        `json:"pap"`
	PapLabel         string     `json:"papLabel"`
	Tags             []string   `json:"tags"`
	Ioc              bool       `json:"ioc"`
	Sighted          bool       `json:"sighted"`
	SightedAt        int64      `json:"sightedAt"`
	Reports          struct{}   `json:"reports"`
	Message          string     `json:"message"`
	ExtraData        ExtraData `json:"extraData"`
	IgnoreSimilarity bool       `json:"ignoreSimilarity"`
}

type ExtraData struct {
	Links *Links `json:"links,omitempty"`
}

func (or *ObservableResponse) UnmarshalJSON(data []byte) error {
	shadow := new(shadowObservableResponse)
	err := json.Unmarshal(data, &shadow)
	if err != nil {
		return err
	}

	or.Id = shadow.Id
	or.Type = shadow.Type
	or.CreatedBy = shadow.CreatedBy
	or.UpdatedBy = shadow.UpdatedBy
	or.CreatedAt = convertInt64ToTime(shadow.CreatedAt)
	or.UpdatedAt = convertInt64ToTime(shadow.UpdatedAt)
	or.DataType = shadow.DataType
	or.Data = shadow.Data
	or.StartDate = convertInt64ToTime(shadow.StartDate)
	or.Attachment = shadow.Attachment
	or.Tlp = shadow.Tlp
	or.TlpLabel = shadow.TlpLabel
	or.Pap = shadow.Pap
	or.PapLabel = shadow.PapLabel
	or.Tags = shadow.Tags
	or.Ioc = shadow.Ioc
	or.Sighted = shadow.Sighted
	or.SightedAt = convertInt64ToTime(shadow.SightedAt)
	or.Reports = shadow.Reports
	or.Message = shadow.Message
	or.ExtraData = shadow.ExtraData
	or.IgnoreSimilarity = shadow.IgnoreSimilarity
	return nil
}

type ObservableTypeResponse struct {
	Id           string    `json:"_id"`
	Type         string    `json:"_type"`
	CreatedAt    time.Time `json:"_createdAt"`
	CreatedBy    string    `json:"_createdBy"`
	Name         string    `json:"name"`
	IsAttachment bool      `json:"isAttachment"`
}

type shadowObservableTypeResponse struct {
	Id           string `json:"_id"`
	Type         string `json:"_type"`
	CreatedAt    int64  `json:"_createdAt"`
	CreatedBy    string `json:"_createdBy"`
	Name         string `json:"name"`
	IsAttachment bool   `json:"isAttachment"`
}

func (or *ObservableTypeResponse) UnmarshalJSON(data []byte) error {
	shadow := new(shadowObservableTypeResponse)
	err := json.Unmarshal(data, &shadow)
	if err != nil {
		return err
	}

	or.Id = shadow.Id
	or.Type = shadow.Type
	or.CreatedAt = convertInt64ToTime(shadow.CreatedAt)
	or.CreatedBy = shadow.CreatedBy
	or.Name = shadow.Name
	or.IsAttachment = shadow.IsAttachment

	return nil
}

// executeObservableTypeSearchQuery is a helper function to do query related searches
func (hive *Hivedata) executeObservableTypeSearchQuery(query []byte) ([]ObservableTypeResponse, error) {
	url, err := url.JoinPath(hive.Url, "/api/v1/query")
	if err != nil {
		return nil, err
	}
	ret, err := hive.webRequest(url, POST, query)

	if err != nil {
		return nil, err
	}

	var parsedRet []ObservableTypeResponse
	err = json.Unmarshal(ret, &parsedRet)
	return parsedRet, err
}

// executeObservableSearchQuery is a helper function to do query related searches
func (hive *Hivedata) executeObservableSearchQuery(query []byte) ([]ObservableResponse, error) {
	url, err := url.JoinPath(hive.Url, "/api/v1/query")
	if err != nil {
		return nil, err
	}
	ret, err := hive.webRequest(url, POST, query)

	if err != nil {
		return nil, err
	}

	var parsedRet []ObservableResponse
	err = json.Unmarshal(ret, &parsedRet)
	return parsedRet, err
}

// GetObservableTypes returns all types an observable can be
func (hive *Hivedata) GetObservableTypes() ([]ObservableTypeResponse, error) {
	query, err := hive.createSearchQuery(
		SearchQuery{Name: "listObservableType"},
		SearchQuery{Name: "sort", Sort: &[1]map[string]string{{"_updatedAt": "desc"}}},
	)
	if err != nil {
		return nil, err
	}

	return hive.executeObservableTypeSearchQuery(query)
}

// AddCaseObservable adds observables to an existing case.
func (hive *Hivedata) AddCaseObservable(incidentNumber int, observable *Observable) error {
	url, err := url.JoinPath(hive.Url, "/api/v1/case/", strconv.Itoa(incidentNumber), "/observable")
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

// AddCaseObservableFile adds a file as an observable to a case.
func (hive *Hivedata) AddCaseObservableFile(incidentNumber int, observable *Observable, file *os.File) error {
	url, err := url.JoinPath(hive.Url, "/api/v1/case/", strconv.Itoa(incidentNumber), "/observable")
	if err != nil {
		return err
	}

	jsondata, err := json.Marshal(observable)
	if err != nil {
		return err
	}

	_, err = hive.webRequestMultiPart(url, POST, jsondata, file)
	return err
}

// GetCaseObservables returns all observables associated with a case
// It returns an observable slice or an error
func (hive *Hivedata) GetCaseObservables(caseId int) ([]ObservableResponse, error) {
	caseNumber := strconv.Itoa(caseId)
	query, err := hive.createSearchQuery(
		SearchQuery{Name: "getCase", IdOrName: caseNumber},
		SearchQuery{Name: "observables"},
	)
	if err != nil {
		return nil, err
	}

	return hive.executeObservableSearchQuery(query)
}

// DeleteObservable deletes an observable
// Only returns data if an error occured
func (hive *Hivedata) DeleteObservable(observableID string) error {
	url, err := url.JoinPath(hive.Url, "/api/v1/observable/", observableID)
	if err != nil {
		return err
	}

	_, err = hive.webRequest(url, DELETE, nil)
	return err
}

// DeleteObservable deletes an observable specified by the ID got through GetObservables
// Only returns data if an error occured
func (hive *Hivedata) UpdateObservable(observableID string, observable *Observable) error {
	url, err := url.JoinPath(hive.Url, "/api/v1/observable/", observableID)
	if err != nil {
		return err
	}

	jsonsearch, err := json.Marshal(observable)
	if err != nil {
		return err
	}

	_, err = hive.webRequest(url, PATCH, jsonsearch)
	return err
}

// Get a single Observable
// It returns a pointer to an observable object or an error
func (hive *Hivedata) GetObservable(observableID string) (*ObservableResponse, error) {
	url, err := url.JoinPath(hive.Url, "/api/v1/observable/", observableID)
	if err != nil {
		return nil, err
	}

	ret, err := hive.webRequest(url, GET, nil)
	if err != nil {
		return nil, err
	}

	var parsedRet *ObservableResponse
	err = json.Unmarshal(ret, &parsedRet)
	return parsedRet, err
}

// GetCaseObservableFiltered returns a single specified observable associated with a case filtered on a field & value
// It returns an observable slice or an error
func (hive *Hivedata) GetCaseObservablesFiltered(caseId int, queryfield, queryvalue string) ([]ObservableResponse, error) {
	caseNumber := strconv.Itoa(caseId)
	query, err := hive.createSearchQuery(
		SearchQuery{Name: "getCase", IdOrName: caseNumber},
		SearchQuery{Name: "observables"},
		SearchQuery{Name: "filter", Eq: &Filter{Field: queryfield, Value: &queryvalue}})
	if err != nil {
		return nil, err
	}

	return hive.executeObservableSearchQuery(query)
}

// Find an observable globally
// It returns a pointer to an ObservableResponse slice or an error
// Be aware that the ObservableResponse will contain a ExtraData field which contains a HiveCaseResponse or HiveAlerResponse object
func (hive *Hivedata) FindObservable(value string) ([]ObservableResponse, error) {
	query, err := hive.createSearchQuery(
		SearchQuery{Name: "listObservable"},
		SearchQuery{Name: "filter", And: &[]Filter{{Field: "keyword", Value: &value}}},
		SearchQuery{Name: "sort", Sort: &[1]map[string]string{{"_createdAt": "desc"}}},
		SearchQuery{Name: "page", ScopeFrom: 0, ScopeTo: 10, ExtraData: []string{"links"}})
	if err != nil {
		return nil, err
	}

	return hive.executeObservableSearchQuery(query)
}
