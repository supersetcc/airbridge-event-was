package webapp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"

	iris "github.com/kataras/iris"
	uuid "github.com/satori/go.uuid"
)

type EventLogKwargs struct {
	AppID         *string     `json:"app_id"`
	AppName       string      `json:"app_name"`
	EventCategory int         `json:"event_category"`
	DeviceUUID    string      `json:"device_uuid"`
	Data          interface{} `json:"data"`
}

type EventLog struct {
	WhatToDo      string         `json:"what_to_do"`
	LogUUID       string         `json:"log_uuid"`
	RecvTimestamp int64          `json:"recv_timestampp"`
	Kwargs        EventLogKwargs `json:"kwargs"`
}

type MobileEvent struct {
	SdkVersion       string `json:"sdkVersion"`
	RequestTimestamp int    `json:"requestTimestamp"`
	EventTimestamp   int    `json:"eventTimestamp"`
	EventUUID        string `json:"eventUUID"`

	ClientData struct {
		OSVersion  string `json:"osVersion"`
		DeviceType string `json:"deviceType"`

		DeferredKey struct {
			DeviceType string `json:"deviceType"`
			OSVersion  string `json:"osVersion"`
		}
	}

	Device struct {
		DeviceModel string `json:"deviceModel"`
		DeviceUUID  string `json:"deviceUUID"`
		OSName      string `json:"osName"`
		OSVersion   string `json:"osVersion"`
	} `json:"device"`

	Browser struct {
		ClientID string `json:"clientID"`
	} `json:"browser"`

	EventData struct {
		TransactionID string `json:"transactionID"`
		ShortID       string `json:"shortID"`

		TrackingData struct {
			Channel string            `json:"channel"`
			Params  map[string]string `json:"params"`
		} `json:"trackingData"`
	} `json:"eventData"`
}

type MobileEventResponse struct {
	ResultMessage string      `json:"resultMessage"`
	Resource      interface{} `json:"resource"`
	At            string      `json:"at"`
}

const (
	EXCEPTION_MSG_VALIDATION    = "Invalid request. Please RTFM. :P"
	EXCEPTION_MSG_GENERAL       = "Sorry. :( Your request was temporarily failed. This issue is reported to our log system. We will fix it as soon."
	EXCEPTION_MSG_UNSUPPORTED   = "Request body is empty or method is not a POST"
	EXCEPTION_MSG_AUTHORIZATION = "There is no Authorization property in the request header"

	APP_NAME              = "app_name"
	EVENT_CATEGORY        = "event_category"
	USER_AGENT            = "User-Agent"
	X_FORWARDED_FOR       = "X-Forwarded-For"
	JOB_NAME_MOBILE_EVENT = "handle_mobile_event"
	AUTHORIZATION         = "Authorization"
)

func GetDeviceModel(event MobileEvent) string {
	if event.Device.DeviceModel != "" {
		return event.Device.DeviceModel
	}

	if event.ClientData.DeviceType != "" {
		return event.ClientData.DeviceType
	}

	if event.ClientData.DeferredKey.DeviceType != "" {
		return event.ClientData.DeferredKey.DeviceType
	}

	return ""
}

func GetOSVersion(event MobileEvent) string {
	if event.Device.OSVersion != "" {
		return event.Device.OSVersion
	}

	if event.ClientData.OSVersion != "" {
		return event.ClientData.OSVersion
	}

	if event.ClientData.DeferredKey.OSVersion != "" {
		return event.ClientData.DeferredKey.OSVersion
	}

	return ""
}

func (app *WebApp) HandleMobileEventReceiver(ic iris.Context) {
	txn, err := app.logging.NewTransaction()
	defer txn.End()

	if err != nil {
		WriteError(ic, 500, EXCEPTION_MSG_GENERAL, "newrelic error")
	}

	request := ic.Request()

	// Authorization Header
	authorization := request.Header.Get(AUTHORIZATION)
	if authorization == "" {
		txn.AddAttribute("http-response-status-code", 401)
		WriteError(ic, 401, EXCEPTION_MSG_AUTHORIZATION, "")
		return
	}

	rawData, err := ioutil.ReadAll(ic.Request().Body)
	if err != nil {
		txn.NoticeError(err)
		txn.AddAttribute("http-response-status-code", 500)
		WriteError(ic, 500, EXCEPTION_MSG_GENERAL, err.Error())
		return
	}

	if len(rawData) == 0 {
		txn.AddAttribute("http-response-status-code", 400)
		WriteError(ic, 400, EXCEPTION_MSG_VALIDATION, "missing body")
		return
	}

	logUUID, err := uuid.NewV4()
	if err != nil {
		txn.NoticeError(err)
		txn.AddAttribute("http-response-status-code", 500)
		WriteError(ic, 500, EXCEPTION_MSG_GENERAL, err.Error())
		return
	}

	// extract IP address from X-Forwared-For
	xForwardedString := request.Header.Get(X_FORWARDED_FOR)
	clientIP := ParseClientIPFromXForwarededFor(xForwardedString)
	txn.AddAttribute("clientIP", clientIP)

	appName := ic.Params().Get(APP_NAME)
	if appName == "" {
		txn.AddAttribute("http-response-status-code", 400)
		WriteError(ic, 400, EXCEPTION_MSG_VALIDATION, "missing 'app_name'")
		return
	}

	txn.AddAttribute("appName", appName)

	eventCategory, err := ic.Params().GetInt(EVENT_CATEGORY)
	if err != nil {
		txn.NoticeError(err)
		txn.AddAttribute("http-response-status-code", 500)
		WriteError(ic, 500, EXCEPTION_MSG_GENERAL, err.Error())
		return
	}

	txn.AddAttribute("eventCategory", eventCategory)

	decoded := map[string]interface{}{}
	if err := json.Unmarshal(rawData, &decoded); err != nil {
		txn.NoticeError(err)
		txn.AddAttribute("http-response-status-code", 500)
		WriteError(ic, 500, EXCEPTION_MSG_GENERAL, err.Error())
		return
	}

	// assign clientIP
	if dm := reflect.ValueOf(decoded["device"]); dm.Kind() == reflect.Map {
		dm.SetMapIndex(reflect.ValueOf("clientIP"), reflect.ValueOf(clientIP))
	}

	mobileEvent := MobileEvent{}
	if err := json.Unmarshal(rawData, &mobileEvent); err != nil {
		txn.NoticeError(err)
		txn.AddAttribute("http-response-status-code", 500)
		WriteError(ic, 500, EXCEPTION_MSG_GENERAL, err.Error())
		return
	}

	payload := EventLog{
		WhatToDo:      JOB_NAME_MOBILE_EVENT, // what_to_do
		LogUUID:       logUUID.String(),      // log_uuid
		RecvTimestamp: CurrentTimestamp(),    // recv_timestamp
		Kwargs: EventLogKwargs{
			AppID:         nil,
			AppName:       appName,
			Data:          decoded, // data
			EventCategory: eventCategory,
			DeviceUUID:    mobileEvent.Device.DeviceUUID,
		},
	}

	// generate kafka partition key
	osVersion := GetOSVersion(mobileEvent)
	deviceModel := GetDeviceModel(mobileEvent)
	appSubdomain := appName
	remoteAddr := clientIP
	pk := GenerateKafkaPartitionKey(osVersion, deviceModel, appSubdomain, remoteAddr)

	encoded, err := json.Marshal(payload)
	if err != nil {
		txn.NoticeError(err)
		txn.AddAttribute("http-response-status-code", 500)
		WriteError(ic, 500, EXCEPTION_MSG_GENERAL, err.Error())
		return
	}

	if err := app.producer.Publish("airbridge-raw-events", pk, encoded); err != nil {
		txn.NoticeError(err)
		txn.AddAttribute("http-response-status-code", 500)
		WriteError(ic, 500, EXCEPTION_MSG_GENERAL, err.Error())
		return
	}

	response := MobileEventResponse{
		ResultMessage: fmt.Sprintf("Event(%d) is successfully proccessed.", eventCategory),
		Resource:      new(map[string]string),
		At:            TimeToStr(KSTNow()),
	}
	WriteResponse(ic, response)

	txn.AddAttribute("http-response-status-code", 200)
}

func (app *WebApp) HandleUnsupportedMethod(ic iris.Context) {
	WriteError(ic, 400, EXCEPTION_MSG_UNSUPPORTED, "")
}
