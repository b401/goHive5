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

// UserResponse contains all available data of a user object on the hive
type UserResponse struct {
	Id                  string            `json:"_id"`
	CreatedBy           string            `json:"_createdBy"`
	UpdatedBy           string            `json:"_updatedBy,omitempty"`
	CreatedAt           time.Time         `json:"_createdAt"`
	UpdatedAt           time.Time         `json:"_updatedAt,omitempty"`
	Login               string            `json:"login"`
	Name                string            `json:"name"`
	Email               string            `json:"email,omitempty"`
	HasKey              bool              `json:"hasKey"`
	HasPassword         bool              `json:"hasPassword"`
	HasMFA              bool              `json:"hasMFA"`
	Locked              bool              `json:"locked"`
	Profile             string            `json:"profile"`
	Permissions         []string          `json:"permissions,omitempty"`
	Organisation        string            `json:"organisation"`
	Avatar              string            `json:"avatar,omitempty"`
	Organisations       []Organisations   `json:"organisations,omitempty"`
	Type                string            `json:"type"`
	DefaultOrganisation string            `json:"defaultOrganisation"`
	ExtraData           map[string]string `json:"extraData"`
}

type Links struct {
	ToOrganisation string `json:"toOrganisation"`
	Avatar         string `json:"avatar,omitempty"`
	LinkType       string `json:"linkType"`
	OtherLinkType  string `json:"otherLinkType"`
}

type Organisations struct {
	OrganisationID string  `json:"organisationId"`
	Organisation   string  `json:"organisation"`
	Profile        string  `json:"profile"`
	Avatar         string  `json:"avatar"`
	Links          []Links `json:"links"`
}

type shadowUserResponse struct {
	Id                  string            `json:"_id"`
	CreatedBy           string            `json:"_createdBy"`
	UpdatedBy           string            `json:"_updatedBy,omitempty"`
	CreatedAt           int64             `json:"_createdAt"`
	UpdatedAt           int64             `json:"_updatedAt,omitempty"`
	Login               string            `json:"login"`
	Name                string            `json:"name"`
	Email               string            `json:"email,omitempty"`
	HasKey              bool              `json:"hasKey"`
	HasPassword         bool              `json:"hasPassword"`
	HasMFA              bool              `json:"hasMFA"`
	Locked              bool              `json:"locked"`
	Profile             string            `json:"profile"`
	Permissions         []string          `json:"permissions,omitempty"`
	Organisation        string            `json:"organisation"`
	Avatar              string            `json:"avatar,omitempty"`
	Organisations       []Organisations   `json:"organisations,omitempty"`
	Type                string            `json:"type"`
	DefaultOrganisation string            `json:"defaultOrganisation"`
	ExtraData           map[string]string `json:"extraData"`
}

func (u *UserResponse) UnmarshalJSON(data []byte) error {
	shadow := new(shadowUserResponse)
	err := json.Unmarshal(data, &shadow)
	if err != nil {
		return err
	}

	u.Id = shadow.Id
	u.CreatedBy = shadow.CreatedBy
	u.UpdatedBy = shadow.UpdatedBy
	u.CreatedAt = convertInt64ToTime(shadow.CreatedAt)
	u.UpdatedAt = convertInt64ToTime(shadow.UpdatedAt)
	u.Login = shadow.Login
	u.Name = shadow.Name
	u.Email = shadow.Email
	u.HasKey = shadow.HasKey
	u.HasPassword = shadow.HasPassword
	u.HasMFA = shadow.HasMFA
	u.Locked = shadow.Locked
	u.Profile = shadow.Profile
	u.Permissions = shadow.Permissions
	u.Organisation = shadow.Organisation
	u.Avatar = shadow.Avatar
	u.Organisations = shadow.Organisations
	u.Type = shadow.Type
	u.DefaultOrganisation = shadow.DefaultOrganisation
	u.ExtraData = shadow.ExtraData

	return nil
}

// executeCommentSearchQuery is a helper function to do query related searches
func (hive *Hivedata) executeUserSearchQuery(query []byte) ([]UserResponse, error) {
	url, err := url.JoinPath(hive.Url, "/api/v1/query")
	if err != nil {
		return nil, err
	}

	ret, err := hive.webRequest(url, POST, query)
	if err != nil {
		return nil, err
	}

	var parsedRet []UserResponse
	err = json.Unmarshal(ret, &parsedRet)
	return parsedRet, err
}

// GetVisibleUsers returns all users that are visible to the service
func (hive *Hivedata) GetVisibleUsers() ([]UserResponse, error) {

	query, err := hive.createSearchQuery(
		SearchQuery{Name: "listVisibleUsers"},
	)
	if err != nil {
		return nil, err
	}

	return hive.executeUserSearchQuery(query)
}
