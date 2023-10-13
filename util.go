/*
thehive5 implements functionality to interact with the most recent version of thehive.
https://www.strangebee.com/thehive/
*/
package thehive5

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Attachment struct {
	Id          string `json:"id"`
	ContentType string `json:"contentType"`
	Name        string `json:"name"`
}

// webRequest is an internal helper to build the right webrequest structure
// it adds additional headers & returns the json body
// Unknown status codes get returned as error
func (hive *Hivedata) webRequest(url string, m method, body []byte) ([]byte, error) {
	b := bytes.NewReader(body)

	req, err := http.NewRequest(string(m), url, b)
	if err != nil {
		return nil, err
	}

	// prepare headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", hive.Apikey))
	resp, err := hive.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check thehive response code and determine if we need to return an error
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errorResp ApiErrorResponse
		err = json.Unmarshal(responseBody, &errorResp)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal api http error")
		}

		return nil, fmt.Errorf("API error: %v", errorResp)

	}

	return responseBody, nil
}

// webRequestMultiPart is an internal helper to build the right webrequest structure for uploads
// it adds additional headers & returns the json body
// Unknown status codes get returned as error
func (hive *Hivedata) webRequestMultiPart(url string, m method, body []byte, file *os.File) ([]byte, error) {
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	partWriter, err := writer.CreateFormField("_json")
	if err != nil {
		return nil, err
	}

	_, err = partWriter.Write(body)
	if err != nil {
		return nil, err
	}

	if file != nil {
		fileWriter, err := writer.CreateFormFile("attachment", filepath.Base(file.Name()))
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(fileWriter, file)
		if err != nil {
			return nil, err
		}
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(string(m), url, &b)
	if err != nil {
		return nil, err
	}

	// prepare headers
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", hive.Apikey))
	resp, err := hive.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check thehive response code and determine if we need to return an error
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errorResp ApiErrorResponse
		err = json.Unmarshal(responseBody, &errorResp)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal api http error")
		}

		return nil, fmt.Errorf("API error: %v", errorResp)

	}

	return responseBody, nil
}

// Helper function to build a search query for the /query endpoint
func (hive *Hivedata) createSearchQuery(filters ...SearchQuery) ([]byte, error) {
	searchquery := HiveSearch{filters}
	return json.Marshal(searchquery)
}

// A GenericResponse can be used for certain API calls to the hive
type ApiErrorResponse struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// Error() reciever implementation to display API errors
func (h ApiErrorResponse) Error() string {
	return fmt.Sprintf("[x] HiveAPI HTTP %s: %s", h.Type, h.Message)
}

// A method is used for HTTP calls to determine what http method should be used
type method string

// Constant to handle HTTP methods
const (
	POST   method = "POST"
	GET    method = "GET"
	PATCH  method = "PATCH"
	DELETE method = "DELETE"
)

// method String returns the http method as a string
func (m method) String() string {
	switch m {
	case POST:
		return "post"
	case GET:
		return "get"
	case PATCH:
		return "patch"
	case DELETE:
		return "delete"
	}
	return "unknown"
}

// Helper function to convert return values into time.Time
func convertInt64ToTime(source int64) time.Time {
	var v time.Time
	if source != 0 {
		return time.UnixMilli(source)
	}
	return v // Return zero time if source is zero
}

// Heler function to convert return values into time.Duration
func convertInt64ToDuration(source int64) time.Duration {
	var v time.Duration
	if source != 0 {
		return time.Duration(source) * time.Millisecond
	}
	return v
}

// A method is used for HTTP calls to determine what http method should be used
type Severity int
type Tlp int
type Pap int

// Constant to handle HTTP methods
const (
	SeverityLow Severity = iota + 1
	SeverityMedium
	SeverityHigh
	SeverityCritical
)

func (s Severity) String() string {
	var value string
	switch s {
	case SeverityLow:
		value = "Low"
	case SeverityMedium:
		value = "Medium"
	case SeverityHigh:
		value = "High"
	case SeverityCritical:
		value = "Critical"
	}

	return fmt.Sprintf("%v", value)
}

func (s *Severity) FromString(v string) error {
	var value Severity
	switch strings.ToLower(v) {
	case "low":
		value = SeverityLow
	case "medium":
		value = SeverityMedium
	case "high":
		value = SeverityHigh
	case "critical":
		value = SeverityCritical
	}

	*s = value

	return nil
}

const (
	TlpClear Tlp = iota
	TlpGreen
	TlpAmber
	TlpAmber_Strict
	TlpRed
)

const (
	PapClear Pap = iota
	PapGreen
	PapAmber
	PapRed
)

func (p Pap) String() string {
	var value string
	switch p {
	case PapClear:
		value = "clear"
	case PapGreen:
		value = "green"
	case PapAmber:
		value = "amber"
	case PapRed:
		value = "red"
	}

	return fmt.Sprintf("%v", value)
}

func (p *Pap) FromString(v string) error {
	var value Pap
	switch strings.ToLower(v) {
	case "clear":
		value = PapClear
	case "green":
		value = PapGreen
	case "amber":
		value = PapAmber
	case "red":
		value = PapRed
	default:
		return fmt.Errorf("unknown PAP value. Allowed: clear,green,amber,red")
	}

	*p = value

	return nil
}

func (t Tlp) String() string {
	var value string
	switch t {
	case TlpClear:
		value = "clear"
	case TlpGreen:
		value = "green"
	case TlpAmber:
		value = "amber"
	case TlpAmber_Strict:
		value = "amber+strict"
	case TlpRed:
		value = "red"
	}

	return fmt.Sprintf("%v", value)
}

func (t *Tlp) FromString(v string) error {
	var value Tlp
	switch strings.ToLower(v) {
	case "clear":
		value = TlpClear
	case "green":
		value = TlpGreen
	case "amber":
		value = TlpAmber
	case "amber+strict":
		value = TlpAmber_Strict
	case "red":
		value = TlpRed
	default:
		return fmt.Errorf("unknown TLP value: %s. Allowed: clear,green,amber,amber+strict,red", v)
	}

	*t = value

	return nil
}

func detectCustomFieldType(s interface{}) (*string, error) {
	var str string
	switch s.(type) {
	// standard tests
	case string:
		str = "string"
	case int:
		str = "integer"
	case bool:
		str = "boolean"
	case float32:
		str = "float"
	case float64:
		str = "float"
	case time.Time:
		str = "date"
	}

	// We check if the interface is a valid url
	if str == "string" {
		_, err := url.ParseRequestURI(s.(string))
		if err == nil {
			str = "url"
		}
	}

	if len(str) != 0 {
		return &str, nil
	}

	return nil, fmt.Errorf("can't find valid customfieldtype for %v", s)

}
