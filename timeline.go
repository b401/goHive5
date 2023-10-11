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
	"time"
)

// TimelineEventResponse contains the response of the timeline api endpoint
type TimelineEventResponse struct {
	Id          string    `json:"_id"`
	Type        string    `json:"_type"`
	CreatedBy   string    `json:"_createdBy"`
	UpdatedBy   string    `json:"_updatedBy,omitempty"`
	CreatedAt   time.Time `json:"_createdAt"`
	UpdatedAt   time.Time `json:"_updatedAt,omitempty"`
	Date        time.Time `json:"date"`
	EndDate     time.Time `json:"endDate,omitempty"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
}

type shadowTimelineEventResponse struct {
	Id          string `json:"_id"`
	Type        string `json:"_type"`
	CreatedBy   string `json:"_createdBy"`
	UpdatedBy   string `json:"_updatedBy,omitempty"`
	CreatedAt   int64  `json:"_createdAt"`
	UpdatedAt   int64  `json:"_updatedAt,omitempty"`
	Date        int64  `json:"date"`
	EndDate     int64  `json:"endDate,omitempty"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

// Unmarshal thehive5 returned values into the HiveHiveAlertResponse structs. Making sure that int64 gets converted into time.Time
func (ter *TimelineEventResponse) UnmarshalJSON(data []byte) error {
	shadow := new(shadowTimelineEventResponse)
	err := json.Unmarshal(data, &shadow)
	if err != nil {
		return err
	}
	ter.Id = shadow.Id
	ter.Type = shadow.Type
	ter.CreatedBy = shadow.CreatedBy
	ter.UpdatedBy = shadow.UpdatedBy
	ter.CreatedAt = convertInt64ToTime(shadow.CreatedAt)
	ter.UpdatedAt = convertInt64ToTime(shadow.UpdatedAt)
	ter.Date = convertInt64ToTime(shadow.Date)
	ter.EndDate = convertInt64ToTime(shadow.EndDate)
	ter.Title = shadow.Title
	ter.Description = shadow.Description

	return nil
}

// FullTimelineResponse gets returned if the full timeline is requested for a case
type FullTimelineResponse struct {
	TimelineDate time.Time   `json:"date"`
	Kind         string      `json:"kind"`
	Entity       string      `json:"entity"`
	EntityID     string      `json:"entityId"`
	Details      EventDetail `json:"details,omitempty"`
}

// shadow marshal struct
type shadowFullTimelineResponse struct {
	TimelineDate int64       `json:"date"`
	Kind         string      `json:"kind"`
	Entity       string      `json:"entity"`
	EntityID     string      `json:"entityId"`
	Details      EventDetail `json:"details,omitempty"`
}

// Unmarshal thehive5 returned values into the FullTimelineResponse structs. Making sure that int64 gets converted into time.Time
func (ftr *FullTimelineResponse) UnmarshalJSON(data []byte) error {
	shadow := new(shadowFullTimelineResponse)
	err := json.Unmarshal(data, &shadow)
	if err != nil {
		return err
	}

	ftr.TimelineDate = convertInt64ToTime(shadow.TimelineDate)
	ftr.Kind = shadow.Kind
	ftr.Entity = shadow.Entity
	ftr.EntityID = shadow.EntityID
	ftr.Details = shadow.Details

	return nil
}

// EventDetail contains either a Task or a CustomEvent
type EventDetail struct {
	Task        *CaseTaskResponse `json:"task,omitempty"`
	CustomEvent *CustomEvent      `json:"customEvent,omitempty"`
}

type CustomEvent struct {
	Id          string    `json:"_id"`
	Type        string    `json:"_type"`
	CreatedBy   string    `json:"_createdBy"`
	UpdatedAt   time.Time `json:"_updatedAt,omitempty"`
	CreatedAt   time.Time `json:"_createdAt"`
	Title       string    `json:"title"`
	Date        time.Time `json:"date"`
	EndDate     time.Time `json:"endDate,omitempty"`
	Description string    `json:"description,omitempty"`
	UpdatedBy   string    `json:"_updatedBy,omitempty"`
}

type shadowCustomEvent struct {
	Id          string `json:"_id"`
	Type        string `json:"_type"`
	CreatedBy   string `json:"_createdBy"`
	UpdatedAt   int64  `json:"_updatedAt,omitempty"`
	CreatedAt   int64  `json:"_createdAt"`
	Title       string `json:"title"`
	Date        int64  `json:"date"`
	EndDate     int64  `json:"endDate,omitempty"`
	Description string `json:"description,omitempty"`
	UpdatedBy   string `json:"_updatedBy,omitempty"`
}

// Unmarshal thehive5 returned values into the Task structs. Making sure that int64 gets converted into time.Time
func (c *CustomEvent) UnmarshalJSON(data []byte) error {
	shadow := new(shadowCustomEvent)
	err := json.Unmarshal(data, &shadow)
	if err != nil {
		return err
	}

	c.Id = shadow.Id
	c.Type = shadow.Type
	c.CreatedBy = shadow.CreatedBy
	c.UpdatedAt = convertInt64ToTime(shadow.UpdatedAt)
	c.CreatedAt = convertInt64ToTime(shadow.CreatedAt)
	c.Title = shadow.Title
	c.Date = convertInt64ToTime(shadow.Date)
	c.EndDate = convertInt64ToTime(shadow.EndDate)
	c.Description = shadow.Description
	c.UpdatedBy = shadow.UpdatedBy
	return nil
}

// TimelineEvent contains the information about a new or changed event for a case
type TimelineEvent struct {
	Date        time.Time `json:"date,omitempty"`
	EndDate     time.Time `json:"endDate,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
}

// Marshalling the TimelineEvent
func (t *TimelineEvent) MarshalJSON() ([]byte, error) {
	type Alias TimelineEvent

	var dateInt64 *int64
	var endDateInt64 *int64
	if !t.Date.IsZero() {
		// We ensure that all data sent to the hive is in UTC format
		v := t.Date.UTC().UnixMilli()
		dateInt64 = &v
	}
	if !t.EndDate.IsZero() {
		// We ensure that all data sent to the hive is in UTC format
		v := t.EndDate.UTC().UnixMilli()
		endDateInt64 = &v
	}
	return json.Marshal(&struct {
		Date    *int64 `json:"date,omitempty"`
		EndDate *int64 `json:"endDate,omitempty"`
		*Alias
	}{
		Date:    dateInt64,
		EndDate: endDateInt64,
		Alias:   (*Alias)(t),
	})
}

// CreateTimelineEvent creates a new CustomEvent in a case
func (hive *Hivedata) CreateTimelineEvent(caseId int, event *TimelineEvent) (*TimelineEventResponse, error) {
	caseNumber := strconv.Itoa(caseId)
	url, err := url.JoinPath(hive.Url, "/api/v1/case/", caseNumber, "/customEvent")
	if err != nil {
		return nil, err
	}

	jsonEvent, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	ret, err := hive.webRequest(url, POST, jsonEvent)
	if err != nil {
		return nil, err
	}
	var parsedRet *TimelineEventResponse
	err = json.Unmarshal(ret, &parsedRet)
	return parsedRet, err
}

// GetTimeline returns all timeline objects from a case. This includes CustomEvents,Tasks and built-in events
func (hive *Hivedata) GetTimeline(caseId int) ([]FullTimelineResponse, error) {
	caseNumber := strconv.Itoa(caseId)
	url, err := url.JoinPath(hive.Url, "/api/v1/case/", caseNumber, "/timeline")
	if err != nil {
		return nil, err
	}

	ret, err := hive.webRequest(url, GET, nil)
	if err != nil {
		return nil, err
	}

	// anonymous struct to help with decoding
	var wrapper struct {
		Events []FullTimelineResponse `json:"events"`
	}

	err = json.Unmarshal(ret, &wrapper)
	if err != nil {
		return nil, err
	}

	// actual struct that holds the events
	return wrapper.Events, nil
}

// GetTimelineEvent returns a single event
// Returns err if no event was found
func (hive *Hivedata) GetTimelineEvent(caseId int, eventId string) (*EventDetail, error) {
	resp, err := hive.GetTimeline(caseId)
	if err != nil {
		return nil, err
	}

	for _, event := range resp {
		if event.EntityID == eventId {
			detailCopy := event.Details
			return &detailCopy, nil
		}
	}

	return nil, fmt.Errorf("no event found with the ID: %s", eventId)
}

// DeleteTimelineEvent deletes a specific event
// Returns err on failure
func (hive *Hivedata) DeleteTimelineEvent(eventId string) error {
	url, err := url.JoinPath(hive.Url, "/api/v1/customEvent/", eventId)
	if err != nil {
		return err
	}

	_, err = hive.webRequest(url, DELETE, nil)
	return err
}

// UpdateTimelineEvent updates a specified event
// returns err on failure
func (hive *Hivedata) UpdateTimelineEvent(eventId string, event *TimelineEvent) error {
	url, err := url.JoinPath(hive.Url, "/api/v1/customEvent/", eventId)
	if err != nil {
		return err
	}

	jsonEvent, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = hive.webRequest(url, PATCH, jsonEvent)
	return err
}
