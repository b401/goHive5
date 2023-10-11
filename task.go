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

// CaseTask contains all task informations
type CaseTask struct {
	Title       string    `json:"title"`
	Group       string    `json:"group,omitempty"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status,omitempty"`
	Flag        bool      `json:"flag,omitempty"`
	StartDate   time.Time `json:"startDate,omitempty"`
	EndDate     time.Time `json:"endDate,omitempty"`
	DueDate     time.Time `json:"dueDate,omitempty"`
	Order       *int      `json:"order,omitempty"`
	Assignee    string    `json:"assignee,omitempty"`
	Mandatory   bool      `json:"mandatory,omitempty"`
}

// Marshalling the CaseTask requests
func (ct *CaseTask) MarshalJSON() ([]byte, error) {
	type Alias CaseTask

	var (
		startDateInt64 *int64
		endDateInt64   *int64
		dueDateInt64   *int64
	)
	if !ct.StartDate.IsZero() {
		// We ensure that all data sent to the hive is in UTC format
		setDate := ct.StartDate.UTC().UnixMilli()
		startDateInt64 = &setDate
	}
	if !ct.EndDate.IsZero() {
		// We ensure that all data sent to the hive is in UTC format
		endDate := ct.EndDate.UTC().UnixMilli()
		endDateInt64 = &endDate
	}
	if !ct.DueDate.IsZero() {
		// We ensure that all data sent to the hive is in UTC format
		dueDate := ct.DueDate.UTC().UnixMilli()
		dueDateInt64 = &dueDate
	}
	return json.Marshal(&struct {
		StartDate *int64 `json:"startDate,omitempty"`
		EndDate   *int64 `json:"endDate,omitempty"`
		DueDate   *int64 `json:"dueDate,omitempty"`
		*Alias
	}{
		StartDate: startDateInt64,
		EndDate:   endDateInt64,
		DueDate:   dueDateInt64,
		Alias:     (*Alias)(ct),
	})
}

// TaskLog contains all task log informations
type TaskLog struct {
	Message           string       `json:"message"`
	StartDate         time.Time    `json:"startDate,omitempty"`
	IncludeInTimeline time.Time    `json:"includeInTimeline,omitempty"`
	Attachments       *interface{} `json:"attachments,omitempty"`
}

// Marshalling the TaskLog requests
func (tl *TaskLog) MarshalJSON() ([]byte, error) {
	type Alias TaskLog
	var (
		startDateInt64         *int64
		includeInTimelineInt64 *int64
	)

	if !tl.StartDate.IsZero() {
		// We ensure that all data sent to the hive is in UTC format
		tmp := tl.StartDate.UTC().UnixMilli()
		startDateInt64 = &tmp
	}
	if !tl.IncludeInTimeline.IsZero() {
		// We ensure that all data sent to the hive is in UTC format
		tmp := tl.IncludeInTimeline.UTC().UnixMilli()
		includeInTimelineInt64 = &tmp
	}
	return json.Marshal(&struct {
		StartDate         *int64 `json:"startDate,omitempty"`
		IncludeInTimeline *int64 `json:"includeInTimeline,omitempty"`
		*Alias
	}{
		StartDate:         startDateInt64,
		IncludeInTimeline: includeInTimelineInt64,
		Alias:             (*Alias)(tl),
	})
}

// TaskLogResponse contains all task log responded
type TaskLogResponse struct {
	Id                string        `json:"_id"`
	Type              string        `json:"_type"`
	CreatedBy         string        `json:"_createdBy"`
	UpdatedBy         string        `json:"_updatedBy"`
	CreatedAt         time.Time     `json:"_createdAt"`
	UpdatedAt         time.Time     `json:"_updatedAt"`
	Message           string        `json:"message"`
	Date              time.Time     `json:"date"`
	Attachments       []interface{} `json:"attachments"`
	Owner             string        `json:"owner"`
	IncludeInTimeline time.Time     `json:"includeInTimeline"`
	ExtraData         struct{}      `json:"extraData"`
}

// TaskLogResponse contains all task log responded
type shadowTaskLogResponse struct {
	Id                string        `json:"_id"`
	Type              string        `json:"_type"`
	CreatedBy         string        `json:"_createdBy"`
	UpdatedBy         string        `json:"_updatedBy"`
	CreatedAt         int64         `json:"_createdAt"`
	UpdatedAt         int64         `json:"_updatedAt"`
	Message           string        `json:"message"`
	Date              int64         `json:"date"`
	Attachments       []interface{} `json:"attachments"`
	Owner             string        `json:"owner"`
	IncludeInTimeline int64         `json:"includeInTimeline"`
	ExtraData         struct{}      `json:"extraData"`
}

// shadow Unmarshalling function for TaskLogResponse
func (tl *TaskLogResponse) UnmarshalJSON(data []byte) error {
	shadow := new(shadowTaskLogResponse)
	err := json.Unmarshal(data, &shadow)
	if err != nil {
		return err
	}

	tl.Id = shadow.Id
	tl.Type = shadow.Type
	tl.CreatedBy = shadow.CreatedBy
	tl.UpdatedBy = shadow.UpdatedBy
	tl.UpdatedAt = convertInt64ToTime(shadow.UpdatedAt)
	tl.Message = shadow.Message
	tl.Date = convertInt64ToTime(shadow.Date)
	tl.Attachments = shadow.Attachments
	tl.Owner = shadow.Owner
	tl.IncludeInTimeline = convertInt64ToTime(shadow.IncludeInTimeline)
	tl.ExtraData = shadow.ExtraData

	return nil
}

// CaseTaskResponse stores the response of a task that was added to a case in The Hive
type CaseTaskResponse struct {
	Id          string    `json:"_id"`
	Type        string    `json:"_type"`
	CreatedBy   string    `json:"_createdBy"`
	UpdatedBy   string    `json:"_updatedBy"`
	CreatedAt   time.Time `json:"_createdAt"`
	UpdatedAt   time.Time `json:"_updatedAt"`
	Title       string    `json:"title"`
	Group       string    `json:"group"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Flag        bool      `json:"flag"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	Assignee    string    `json:"assignee"`
	Order       int       `json:"order"`
	DueDate     time.Time `json:"dueDate"`
	Mandatory   bool      `json:"mandatory"`
	ExtraData   struct{}  `json:"extraData"`
}

// CaseTaskResponse stores the response of a task that was added to a case in The Hive
type shadowCaseTaskResponse struct {
	Id          string   `json:"_id"`
	Type        string   `json:"_type"`
	CreatedBy   string   `json:"_createdBy"`
	UpdatedBy   string   `json:"_updatedBy"`
	CreatedAt   int64    `json:"_createdAt"`
	UpdatedAt   int64    `json:"_updatedAt"`
	Title       string   `json:"title"`
	Group       string   `json:"group"`
	Description string   `json:"description"`
	Status      string   `json:"status"`
	Flag        bool     `json:"flag"`
	StartDate   int64    `json:"startDate"`
	EndDate     int64    `json:"endDate"`
	Assignee    string   `json:"assignee"`
	Order       int      `json:"order"`
	DueDate     int64    `json:"dueDate"`
	Mandatory   bool     `json:"mandatory"`
	ExtraData   struct{} `json:"extraData"`
}

// shadow Unmarshalling function for CaseTaskResponse
func (ct *CaseTaskResponse) UnmarshalJSON(data []byte) error {
	shadow := new(shadowCaseTaskResponse)
	err := json.Unmarshal(data, &shadow)
	if err != nil {
		return err
	}

	ct.Id = shadow.Id
	ct.Type = shadow.Type
	ct.CreatedBy = shadow.CreatedBy
	ct.UpdatedBy = shadow.UpdatedBy
	ct.CreatedAt = convertInt64ToTime(shadow.CreatedAt)
	ct.UpdatedAt = convertInt64ToTime(shadow.UpdatedAt)
	ct.Title = shadow.Title
	ct.Group = shadow.Group
	ct.Description = shadow.Description
	ct.Status = shadow.Status
	ct.Flag = shadow.Flag
	ct.Mandatory = shadow.Mandatory
	ct.StartDate = convertInt64ToTime(shadow.StartDate)
	ct.EndDate = convertInt64ToTime(shadow.EndDate)
	ct.Assignee = shadow.Assignee
	ct.Order = shadow.Order
	ct.DueDate = convertInt64ToTime(shadow.DueDate)
	ct.ExtraData = shadow.ExtraData

	return nil
}

// executeTaskSearchQuery is a helper function to do query related searches
func (hive *Hivedata) executeTaskSearchQuery(query []byte) ([]CaseTaskResponse, error) {
	url, err := url.JoinPath(hive.Url, "/api/v1/query")
	if err != nil {
		return nil, err
	}

	ret, err := hive.webRequest(url, POST, query)

	if err != nil {
		return nil, err
	}

	var parsedRet []CaseTaskResponse
	err = json.Unmarshal(ret, &parsedRet)
	return parsedRet, err
}

// executeTaskLogSearchQuery is a helper function to do query related searches
func (hive *Hivedata) executeTaskLogSearchQuery(query []byte) ([]TaskLogResponse, error) {
	url, err := url.JoinPath(hive.Url, "/api/v1/query")
	if err != nil {
		return nil, err
	}

	ret, err := hive.webRequest(url, POST, query)

	if err != nil {
		return nil, err
	}

	var parsedRet []TaskLogResponse
	err = json.Unmarshal(ret, &parsedRet)
	return parsedRet, err
}

// UpdateTask updates an existing task
// Only returns data if an error occured
func (hive *Hivedata) UpdateTask(taskId string, task *CaseTask) error {
	url, err := url.JoinPath(hive.Url, "/api/v1/task/", taskId)
	if err != nil {
		return err
	}

	jsondata, err := json.Marshal(task)
	if err != nil {
		return err
	}

	_, err = hive.webRequest(url, PATCH, jsondata)
	return err
}

// AddTaskToCase creates a new task and adds to an existing case
func (hive *Hivedata) AddTaskToCase(caseId int, task *CaseTask) (*CaseTaskResponse, error) {
	url, err := url.JoinPath(hive.Url, "/api/v1/case/", strconv.Itoa(caseId), "/task")
	if err != nil {
		return nil, err
	}

	jsondata, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	ret, err := hive.webRequest(url, POST, jsondata)
	if err != nil {
		return nil, err
	}

	parsedRet := CaseTaskResponse{}
	err = json.Unmarshal(ret, &parsedRet)
	if err != nil {
		return nil, err
	}
	return &parsedRet, err
}

// DeleteTask deletes an existing task from a case
// only returns data if an error occured
func (hive *Hivedata) DeleteTask(taskId string) error {
	url, err := url.JoinPath(hive.Url, "/api/v1/task/", taskId)
	if err != nil {
		return err
	}

	_, err = hive.webRequest(url, DELETE, nil)
	return err
}

// GetTask returns a single CaseTaskResponse object or error
// A task ID must be provided
func (hive *Hivedata) GetTask(taskId string) (*CaseTaskResponse, error) {
	url, err := url.JoinPath(hive.Url, "/api/v1/task/", taskId)
	if err != nil {
		return nil, err
	}

	ret, err := hive.webRequest(url, GET, nil)
	if err != nil {
		return nil, err
	}

	parsedRet := CaseTaskResponse{}
	err = json.Unmarshal(ret, &parsedRet)
	if err != nil {
		return nil, err
	}
	return &parsedRet, nil
}

// GetTaskLog returns all log entries of a task
// A task ID must be provided
func (hive *Hivedata) CreateTaskLog(taskId string, log *TaskLog) (*TaskLogResponse, error) {
	url, err := url.JoinPath(hive.Url, "/api/v1/task/", taskId, "/log")
	if err != nil {
		return nil, err
	}

	jsondata, err := json.Marshal(log)
	if err != nil {
		return nil, err
	}

	ret, err := hive.webRequest(url, POST, jsondata)
	if err != nil {
		return nil, err
	}

	parsedRet := TaskLogResponse{}
	err = json.Unmarshal(ret, &parsedRet)
	if err != nil {
		return nil, err
	}
	return &parsedRet, nil
}

// GetTaskLogs returns all logs associated with a task
// It returns a task log slice or an error
func (hive *Hivedata) GetTaskLogs(taskId string) ([]TaskLogResponse, error) {
	query, err := hive.createSearchQuery(
		SearchQuery{Name: "getTask", IdOrName: &taskId},
		SearchQuery{Name: "logs"},
		SearchQuery{Name: "sort", Sort: &[1]map[string]string{{"_createdAt": "asc"}}})
	if err != nil {
		return nil, err
	}

	return hive.executeTaskLogSearchQuery(query)
}

// GetCaseTasks returns all tasks associated with a case
// It returns a task slice or an error
func (hive *Hivedata) GetCaseTasks(caseId int) ([]CaseTaskResponse, error) {
	caseNumber := strconv.Itoa(caseId)
	query, err := hive.createSearchQuery(
		SearchQuery{Name: "getCase", IdOrName: &caseNumber},
		SearchQuery{Name: "tasks"},
		SearchQuery{Name: "sort", Sort: &[1]map[string]string{{"_createdAt": "asc"}}})
	if err != nil {
		return nil, err
	}

	return hive.executeTaskSearchQuery(query)
}
