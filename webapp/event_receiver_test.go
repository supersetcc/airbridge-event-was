package webapp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/kataras/iris/httptest"
)

const (
	TestRequestDataPath = "../res/test/mobile_event_request_data.txt"
)

type MockEventReceiverLog struct {
	WhatToDo string                 `json:"what_to_do"`
	Kwargs   map[string]interface{} `json:"kwargs"`
}

func readAllMockRequestPayload() []string {
	res := []string{}

	file, err := os.Open(TestRequestDataPath)
	defer file.Close()

	if err != nil {
		return res
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}

	return res
}

func TestMockMobileEventRequestWithoutAuthorized(t *testing.T) {
	expect, _ := NewWebAppExpect(t)
	uri := fmt.Sprintf("/api/v2/apps/%s/events/mobile-app/%d", "ablog", 9162)
	expect.POST(uri).Expect().Status(httptest.StatusUnauthorized)
}

func TestMockWebAppEventRequestWithoutAuthorized(t *testing.T) {
	payloads := readAllMockRequestPayload()
	if len(payloads) == 0 {
		t.Fatalf("could not read test data(%s)", TestRequestDataPath)
	}

	expect, _ := NewWebAppExpect(t)
	uri := fmt.Sprintf("/api/v2/apps/%s/events/mobile-webapp/%d", "ablog", 9160)
	expect.POST(uri).WithText(payloads[0]).Expect().Status(httptest.StatusOK)
}

func TestMockMobileEventRequestWithAuthorized(t *testing.T) {
	expect, _ := NewWebAppExpect(t)
	uri := fmt.Sprintf("/api/v2/apps/%s/events/mobile-app/%d", "ablog", 9162)
	expect.POST(uri).WithHeader("Authorization", "random-authorized-string").Expect().Status(httptest.StatusBadRequest)
}

func TestMockMobileEventRequestBasic(t *testing.T) {
	payloads := readAllMockRequestPayload()
	if len(payloads) == 0 {
		t.Fatalf("could not read test data(%s)", TestRequestDataPath)
	}

	for _, payload := range payloads {
		expect, mp := NewWebAppExpect(t)
		uri := fmt.Sprintf("/api/v2/apps/%s/events/mobile-app/%d", "ablog", 9162)

		request := expect.POST(uri).WithHeader("Authorization", "random-authorized-string")
		request.WithText(payload)
		request.Expect().Status(httptest.StatusOK)

		if mp.LastPublishedTopic != "airbridge-raw-events" {
			t.Fatalf("publish is not sent to 'airbrdige-raw-events' but %v", mp.LastPublishedTopic)
		}
	}
}

func TestAppIDMustGetNullValue(t *testing.T) {
	payloads := readAllMockRequestPayload()
	for _, payload := range payloads {

		expect, mp := NewWebAppExpect(t)
		uri := fmt.Sprintf("/api/v2/apps/%s/events/mobile-app/%d", "ablog", 9162)

		request := expect.POST(uri).WithHeader("Authorization", "random-authorized-string")
		request.WithText(payload)
		request.Expect().Status(httptest.StatusOK)

		var log MockEventReceiverLog
		err := json.Unmarshal(mp.LastPublishedPayload, &log)
		if err != nil {
			t.Fatalf("could not parse queueing message: %v", err)
		}

		fmt.Println(string(mp.LastPublishedPayload))
		if log.Kwargs["app_id"] != nil {
			t.Fatalf("kwargs['app_id'] must have null value")
		}
	}
}

func TestCheckDataResponseFormat(t *testing.T) {
}

func TestReqestPayloadLessThan512Bytes(t *testing.T) {
	uri := fmt.Sprintf("/api/v2/apps/%s/events/mobile-app/%d", "ablog", 9162)
	expect, _ := NewWebAppExpect(t)
	request := expect.POST(uri).WithHeader("Authorization", "random-auth-data")
	request.WithText("{}")
	request.Expect().Status(httptest.StatusOK)
}
