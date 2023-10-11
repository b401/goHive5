/*
thehive5 implements functionality to interact with the most recent version of thehive.
https://www.strangebee.com/thehive/
*/
package thehive5

import (
	"encoding/json"
	"net/url"
	"time"
)

// A CaseTemplate contains the mapping for the thehive5 api
type CaseTemplate struct {
	Name         string         `json:"name"`
	DisplayName  string         `json:"displayName,omitempty"`
	TitlePrefix  string         `json:"titlePrefix,omitempty"`
	Description  string         `json:"description,omitempty"`
	Severity     *Severity      `json:"severity,omitempty"`
	Tags         *[]string      `json:"tags,omitempty"`
	Flag         *bool          `json:"flag,omitempty"`
	Tlp          string         `json:"tlp,omitempty"`
	Pap          string         `json:"pap,omitempty"`
	Summary      string         `json:"summary,omitempty"`
	CustomFields *[]CustomField `json:"customFields"`
}

// Marshalling the alert requests
func (c *CaseTemplate) MarshalJSON() ([]byte, error) {
	type Alias CaseTemplate
	var (
		tlpInt int
		papInt int
	)

	if len(c.Tlp) != 0 {
		var tlp Tlp
		err := tlp.FromString(c.Tlp)
		if err != nil {
			return nil, err
		}
		tlpInt = int(tlp)
	}
	if len(c.Pap) != 0 {
		var pap Pap
		err := pap.FromString(c.Pap)
		if err != nil {
			return nil, err
		}
		papInt = int(pap)
	}

	return json.Marshal(&struct {
		Tlp int `json:"tlp,omitempty"`
		Pap int `json:"pap,omitempty"`
		*Alias
	}{
		Tlp:   tlpInt,
		Pap:   papInt,
		Alias: (*Alias)(c),
	})
}

// CaseTemplateResponse contain the response of thehive5 templates endpoint
type CaseTemplateResponse struct {
	Id            string        `json:"_id"`
	Type          string        `json:"_type"`
	CreatedBy     string        `json:"_createdBy"`
	UpdatedBy     string        `json:"_updatedBy,omitempty"`
	CreatedAt     time.Time     `json:"_createdAt"`
	UpdatedAt     time.Time     `json:"_updatedAt,omitempty"`
	Name          string        `json:"name"`
	DisplayName   string        `json:"displayName"`
	TitlePrefix   string        `json:"titlePrefix,omitempty"`
	Description   string        `json:"description,omitempty"`
	Severity      Severity      `json:"severity,omitempty"`
	SeverityLabel string        `json:"severityLabel,omitempty"`
	Tags          []string      `json:"tags,omitempty"`
	Flag          bool          `json:"flag"`
	Tlp           Tlp           `json:"tlp,omitempty"`
	TlpLabel      string        `json:"tlpLabel,omitempty"`
	Pap           Pap           `json:"pap,omitempty"`
	PapLabel      string        `json:"papLabel,omitempty"`
	Summary       string        `json:"summary,omitempty"`
	CustomFields  []CustomField `json:"customFields,omitempty"`
	Tasks         []CaseTask    `json:"tasks,omitempty"`
	ExtraData     struct{}      `json:"extraData,omitempty"`
}

// CaseTemplateResponse contain the response of thehive5 templates endpoint
type shadowCaseTemplateResponse struct {
	Id            string        `json:"_id"`
	Type          string        `json:"_type"`
	CreatedBy     string        `json:"_createdBy"`
	UpdatedBy     string        `json:"_updatedBy,omitempty"`
	CreatedAt     int64         `json:"_createdAt"`
	UpdatedAt     int64         `json:"_updatedAt,omitempty"`
	Name          string        `json:"name"`
	DisplayName   string        `json:"displayName"`
	TitlePrefix   string        `json:"titlePrefix,omitempty"`
	Description   string        `json:"description,omitempty"`
	Severity      Severity      `json:"severity,omitempty"`
	SeverityLabel string        `json:"severityLabel,omitempty"`
	Tags          []string      `json:"tags,omitempty"`
	Flag          bool          `json:"flag"`
	Tlp           Tlp           `json:"tlp,omitempty"`
	TlpLabel      string        `json:"tlpLabel,omitempty"`
	Pap           Pap           `json:"pap,omitempty"`
	PapLabel      string        `json:"papLabel,omitempty"`
	Summary       string        `json:"summary,omitempty"`
	CustomFields  []CustomField `json:"customFields,omitempty"`
	Tasks         []CaseTask    `json:"tasks,omitempty"`
	ExtraData     struct{}      `json:"extraData,omitempty"`
}

// shadow unmarshal function for CaseTemplate
func (ctr *CaseTemplateResponse) UnmarshalJSON(data []byte) error {
	shadow := new(shadowCaseTemplateResponse)
	err := json.Unmarshal(data, &shadow)
	if err != nil {
		return err
	}

	ctr.Id = shadow.Id
	ctr.Type = shadow.Type
	ctr.CreatedBy = shadow.CreatedBy
	ctr.UpdatedBy = shadow.UpdatedBy
	ctr.CreatedAt = convertInt64ToTime(shadow.CreatedAt)
	ctr.UpdatedAt = convertInt64ToTime(shadow.UpdatedAt)
	ctr.Name = shadow.Name
	ctr.DisplayName = shadow.DisplayName
	ctr.TitlePrefix = shadow.TitlePrefix
	ctr.Description = shadow.Description
	ctr.Severity = shadow.Severity
	ctr.Tags = shadow.Tags
	ctr.Flag = shadow.Flag
	ctr.Tlp = shadow.Tlp
	ctr.Pap = shadow.Pap
	ctr.Summary = shadow.Summary
	ctr.CustomFields = shadow.CustomFields

	return nil
}

// GetCaseTemplate looks up a specific template on thehive5 instance.
// It returns the CaseTemplateResponse of an error on failure
func (hive *Hivedata) GetCaseTemplate(templateName string) (*CaseTemplateResponse, error) {
	url, err := url.JoinPath(hive.Url, "/api/v1/caseTemplate/", templateName)
	if err != nil {
		return nil, err
	}

	ret, err := hive.webRequest(url, GET, nil)
	if err != nil {
		return nil, err
	}

	parsedRet := new(CaseTemplateResponse)
	err = json.Unmarshal(ret, parsedRet)
	return parsedRet, err
}

// DeleteCaseTemplate allows the deletion of templates on thehive5 instance.
// it returns an error only if the deletion failed.
func (hive *Hivedata) DeleteCaseTemplate(templateName string) error {
	url, err := url.JoinPath(hive.Url, "/api/v1/caseTemplate/", templateName)
	if err != nil {
		return err
	}

	_, err = hive.webRequest(url, DELETE, nil)
	return err
}

// UpdateCaseTemplate updates an already existing template on thehive5.
// Submit the template name as first argument and a CaseTemplate object with the attributes you want to overwrite as the second.
func (hive *Hivedata) UpdateCaseTemplate(templateName string, updatedTemplate CaseTemplate) error {
	url, err := url.JoinPath(hive.Url, "/api/v1/caseTemplate/", templateName)
	if err != nil {
		return err
	}

	jsonrequest, err := json.Marshal(updatedTemplate)
	if err != nil {
		return err
	}

	_, err = hive.webRequest(url, PATCH, jsonrequest)
	return err
}
