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

// comment contains informations about a comment object on an alert or case
type CommentResponse struct {
	Id        string            `json:"_id"`
	Type      string            `json:"_type"`
	CreatedBy string            `json:"createdBy"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
	UpdatedBy string            `json:"updatedBy"`
	Message   string            `json:"message"`
	IsEdited  bool              `json:"isEdited"`
	ExtraData map[string]string `json:"extraData"`
}

type Comment struct {
	Message string `json:"message"`
}

// Marshalling the comment converting the time objects into Unixmilli int64 values
func (c *CommentResponse) MarshalJSON() ([]byte, error) {
	type Alias CommentResponse

	var (
		createdAt *int64
	)
	if !c.CreatedAt.IsZero() {
		// We ensure that all data sent to the hive is in UTC format
		t := c.CreatedAt.UTC().UnixMilli()
		createdAt = &t
	}
	return json.Marshal(&struct {
		CreatedAt *int64 `json:"createdAt,omitempty"`
		*Alias
	}{
		CreatedAt: createdAt,
		Alias:     (*Alias)(c),
	})
}

// shadowComment is used to parse the int64 return value into time.Time
type shadowCommentResponse struct {
	Id        string            `json:"_id"`
	Type      string            `json:"_type"`
	CreatedBy string            `json:"createdBy"`
	CreatedAt int64             `json:"createdAt"`
	UpdatedAt int64             `json:"updatedAt"`
	UpdatedBy string            `json:"updatedBy"`
	Message   string            `json:"message"`
	IsEdited  bool              `json:"isEdited"`
	ExtraData map[string]string `json:"extraData"`
}

// Custom unmarshaller for Comment
func (c *CommentResponse) UnmarshalJSON(data []byte) error {
	shadow := new(shadowCommentResponse)
	err := json.Unmarshal(data, &shadow)
	if err != nil {
		return err
	}

	c.Id = shadow.Id
	c.Type = shadow.Type
	c.CreatedBy = shadow.CreatedBy
	c.CreatedAt = convertInt64ToTime(shadow.CreatedAt)
	c.UpdatedAt = convertInt64ToTime(shadow.UpdatedAt)
	c.UpdatedBy = shadow.UpdatedBy
	c.Message = shadow.Message
	c.IsEdited = shadow.IsEdited
	c.ExtraData = shadow.ExtraData
	return nil
}

// executeCommentSearchQuery is a helper function to do query related searches
func (hive *Hivedata) executeCommentSearchQuery(query []byte) ([]CommentResponse, error) {
	url, err := url.JoinPath(hive.Url, "/api/v1/query")
	if err != nil {
		return nil, err
	}

	ret, err := hive.webRequest(url, POST, query)
	if err != nil {
		return nil, err
	}

	var parsedRet []CommentResponse
	err = json.Unmarshal(ret, &parsedRet)
	return parsedRet, err
}

// GetAlertComments returns all comments associated with an alert
// It returns a comment slice or an error
func (hive *Hivedata) GetAlertComments(alertId string) ([]CommentResponse, error) {
	query, err := hive.createSearchQuery(
		SearchQuery{Name: "getAlert", IdOrName: &alertId},
		SearchQuery{Name: "comments"},
	)
	if err != nil {
		return nil, err
	}

	return hive.executeCommentSearchQuery(query)
}

// AddAlertComment adds a comment to an existing alert
// Returns the created comment as Comment or an error
func (hive *Hivedata) AddAlertComment(alertId string, comment *Comment) (*CommentResponse, error) {
	url, err := url.JoinPath(hive.Url, "/api/v1/alert/", alertId, "/comment")
	if err != nil {
		return nil, err
	}

	jsonValue, err := json.Marshal(comment)
	if err != nil {
		return nil, err
	}

	ret, err := hive.webRequest(url, POST, jsonValue)
	if err != nil {
		return nil, err
	}

	var parsedRet CommentResponse
	err = json.Unmarshal(ret, &parsedRet)
	return &parsedRet, err
}

// AddCaseComment adds a comment to an existing case
// Returns the created comment as Comment or an error
func (hive *Hivedata) AddCaseComment(caseId int, comment *Comment) (*CommentResponse, error) {
	caseNumber := strconv.Itoa(caseId)
	url, err := url.JoinPath(hive.Url, "/api/v1/case/", caseNumber, "/comment")
	if err != nil {
		return nil, err
	}

	jsonValue, err := json.Marshal(comment)
	if err != nil {
		return nil, err
	}

	ret, err := hive.webRequest(url, POST, jsonValue)
	if err != nil {
		return nil, err
	}

	var parsedRet CommentResponse
	err = json.Unmarshal(ret, &parsedRet)
	return &parsedRet, err
}

// GetCaseComments returns all comments associated with a case
// It returns a comment slice or an error
func (hive *Hivedata) GetCaseComments(caseId int) ([]CommentResponse, error) {
	caseNumber := strconv.Itoa(caseId)
	query, err := hive.createSearchQuery(
		SearchQuery{Name: "getCase", IdOrName: &caseNumber},
		SearchQuery{Name: "comments"},
	)
	if err != nil {
		return nil, err
	}

	return hive.executeCommentSearchQuery(query)
}

// GetCaseCommentsTimed takes a time object to look back and return all comments for a specific case in that frame.
func (hive *Hivedata) GetCaseCommentsTimed(caseId int, timeframe time.Time) ([]CommentResponse, error) {
	caseNumber := strconv.Itoa(caseId)

	query, err := hive.createSearchQuery(
		SearchQuery{Name: "getCase", IdOrName: &caseNumber},
		SearchQuery{Name: "comments"},
		SearchQuery{Name: "filter", Gte: &Filter{Field: "_createdAt", Value: timeframe.UnixMilli()}},
	)
	if err != nil {
		return nil, err
	}

	return hive.executeCommentSearchQuery(query)
}
