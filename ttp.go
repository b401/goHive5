/*
thehive5 implements functionality to interact with the most recent version of thehive.
https://www.strangebee.com/thehive/
*/
package thehive5

import (
	"encoding/json"
	"net/url"
	"strconv"
	"time"
)

// A Procedure contains TTPs
type Procedure struct {
	PatternId   string    `json:"patternId"`
	OccurDate   time.Time `json:"occurDate"`
	Tactic      *string   `json:"tactic,omitempty"`
	Description *string   `json:"description,omitempty"`
}

// Marshalling the alert requests
func (p *Procedure) MarshalJSON() ([]byte, error) {
	type Alias Procedure

	var occurdateInt64 int64
	if !p.OccurDate.IsZero() {
		// We ensure that all data sent to the hive is in UTC format
		occurdateInt64 = p.OccurDate.UTC().UnixMilli()
	}
	return json.Marshal(&struct {
		OccurDate *int64 `json:"occurDate,omitempty"`
		*Alias
	}{
		OccurDate: &occurdateInt64,
		Alias:     (*Alias)(p),
	})
}

// ProcedureResponse contains the values of the procedure/ttp operations
type ProcedureResponse struct {
	Id          string            `json:"_id,"`
	CreatedAt   time.Time         `json:"_createdAt"`
	CreatedBy   string            `json:"_createdBy"`
	UpdatedAt   time.Time         `json:"_updatedAt,omitempty"`
	UpdatedBy   string            `json:"_updatedBy,omitempty"`
	Description string            `json:"description,omitempty"`
	OccurDate   time.Time         `json:"occurDate"`
	PatternID   string            `json:"patternId,omitempty"`
	PatternName string            `json:"patternName,omitempty"`
	Tactic      string            `json:"tactic"`
	TacticLabel string            `json:"tacticLabel"`
	ExtraData   map[string]string `json:"extraData"`
}

// ProcedureResponse contains the values of the procedure/ttp operations
type shadowProcedureResponse struct {
	Id          string            `json:"_id,"`
	CreatedAt   int64             `json:"_createdAt"`
	CreatedBy   string            `json:"_createdBy"`
	UpdatedAt   int64             `json:"_updatedAt,omitempty"`
	UpdatedBy   string            `json:"_updatedBy,omitempty"`
	Description string            `json:"description,omitempty"`
	OccurDate   int64             `json:"occurDate"`
	PatternID   string            `json:"patternId,omitempty"`
	PatternName string            `json:"patternName,omitempty"`
	Tactic      string            `json:"tactic"`
	TacticLabel string            `json:"tacticLabel"`
	ExtraData   map[string]string `json:"extraData"`
}

// shadow unmarshalling for ProcedureResponse
func (p *ProcedureResponse) UnmarshalJSON(data []byte) error {
	shadow := new(shadowProcedureResponse)
	err := json.Unmarshal(data, &shadow)
	if err != nil {
		return err
	}

	p.Id = shadow.Id
	p.CreatedAt = convertInt64ToTime(shadow.CreatedAt)
	p.CreatedBy = shadow.CreatedBy
	p.UpdatedAt = convertInt64ToTime(shadow.UpdatedAt)
	p.UpdatedBy = shadow.UpdatedBy
	p.Description = shadow.Description
	p.OccurDate = convertInt64ToTime(shadow.OccurDate)
	p.PatternID = shadow.PatternID
	p.PatternName = shadow.PatternName
	p.Tactic = shadow.Tactic
	p.TacticLabel = shadow.TacticLabel
	p.ExtraData = shadow.ExtraData

	return nil
}

// AddAlertProcedure adds a procedure to an existing alert
func (hive *Hivedata) AddAlertProcedure(alertId string, procedure *Procedure) (*ProcedureResponse, error) {
	url, err := url.JoinPath(hive.Url, "/api/v1/alert", alertId, "/procedurej")
	if err != nil {
		return nil, err
	}

	jsonsearch, err := json.Marshal(procedure)
	if err != nil {
		return nil, err
	}

	ret, err := hive.webRequest(url, POST, jsonsearch)
	if err != nil {
		return nil, err
	}

	var parsedRet *ProcedureResponse
	err = json.Unmarshal(ret, parsedRet)
	return parsedRet, err
}

// AddCaseProcedure adds a procedure to an existing case
func (hive *Hivedata) AddCaseProcedure(caseId int, procedure *Procedure) (*ProcedureResponse, error) {
	caseNumber := strconv.Itoa(caseId)
	url, err := url.JoinPath(hive.Url, "/api/v1/case/", caseNumber, "/procedurej")
	if err != nil {
		return nil, err
	}

	jsonsearch, err := json.Marshal(procedure)
	if err != nil {
		return nil, err
	}

	ret, err := hive.webRequest(url, POST, jsonsearch)
	if err != nil {
		return nil, err
	}

	var parsedRet *ProcedureResponse
	err = json.Unmarshal(ret, parsedRet)
	return parsedRet, err
}
