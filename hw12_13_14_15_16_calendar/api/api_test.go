package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	api "github.com/Calendar/hw12_13_14_15_calendar/api/gen/go"
)

var Id int64 = 1234567

func TestAPI(t *testing.T) {
	t.Run("TestAPI", func(t *testing.T) {
		APIPostRequest(t)
		APIPutRequest(t, Id)
		APIDeleteRequest(t, Id)
	})
}

func APIPostRequest(t *testing.T) {
	event, err := APIGetRequest(nil)
	if err != nil {
		t.Fatal(err)
	}
	requestBody := api.Event{
		Id:             Id,
		Title:          "test api requests in unit test",
		DataStart:      time.Now(),
		DataEnd:        time.Now().Add(60 * time.Minute),
		Description:    "test api requests in unit test: create, delete, update, read.",
		UserId:         1,
		TimeUntilEvent: 15,
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal(err)
	}
	//requestBody, _ := json.Marshal(map[string]string{
	//	"ID":             "2",
	//	"Title":          "test api requests in unit test",
	//	"DataStart":      "2025-06-13T21:36:00Z",
	//	"DataEnd":        "2025-06-13T22:36:00Z",
	//	"Description":    "test api requests in unit test: create, delete, update, read.",
	//	"IdUser":         "1",
	//	"TimeUntilEvent": "15",
	//})
	resp, err := http.Post("http://localhost:8080/events", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatal(fmt.Errorf("status code is %d", resp.StatusCode))
	}
	event = append(event, requestBody)
	_, err = APIGetRequest(event)
	if err != nil {
		APIDeleteRequest(t, Id)
		t.Fatal(err)
	}

}

func APIGetRequest(event []api.Event) ([]api.Event, error) {
	resp, err := http.Get("http://localhost:8080/events")
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Expected status OK, got  %v", resp.StatusCode)
	}
	events := make([]api.Event, 0)
	err = json.NewDecoder(resp.Body).Decode(&events)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	if event == nil {
		return events, nil
	}
	if len(event) != len(events) {
		return nil, errors.New("the response length does not match")
	}
	return events, nil
}
func APIDeleteRequest(t *testing.T, id int64) {
	event, err := APIGetRequest(nil)
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/event/"+strconv.FormatInt(id, 10), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatal(fmt.Errorf("status code is %d", resp.StatusCode))
	}
	_, err = APIGetRequest(event[:len(event)-1])
	if err != nil {
		t.Fatal(err)
	}
}
func APIPutRequest(t *testing.T, id int64) {
	event, err := APIGetRequest(nil)
	if err != nil {
		t.Fatal(err)
	}
	requestBody := api.Event{
		Id:             id,
		Title:          "test api ",
		DataStart:      time.Now(),
		DataEnd:        time.Now().Add(120 * time.Minute),
		Description:    "test api requests in unit test: create, delete, update, read.",
		UserId:         1,
		TimeUntilEvent: 45,
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/event/"+strconv.FormatInt(id, 10), bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatal(fmt.Errorf("status code is %d", resp.StatusCode))
	}
	_, err = APIGetRequest(event)
	if err != nil {
		t.Fatal(err)
	}
}
