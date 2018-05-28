package webapp

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"reflect"

	iris "github.com/kataras/iris"
	newrelic "github.com/newrelic/go-agent"
	uuid "github.com/satori/go.uuid"
)

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

func (app *WebApp) HandleEventReceiverMobile(ic iris.Context) {
	txn, _ := app.logging.NewTransaction()
	defer txn.End()

	// Authorization Header
	request := ic.Request()
	authorization := request.Header.Get(AUTHORIZATION)
	if authorization == "" {
		txn.AddAttribute("http-response-status-code", 401)
		WriteError(ic, 401, EXCEPTION_MSG_AUTHORIZATION, "")
		return
	}

	app.handleEvent(ic, txn)
}

func (app *WebApp) HandleEventReceiverWebApp(ic iris.Context) {
	txn, _ := app.logging.NewTransaction()
	defer txn.End()

	app.handleEvent(ic, txn)
}

func (app *WebApp) HandleUnsupportedMethod(ic iris.Context) {
	WriteError(ic, 400, EXCEPTION_MSG_UNSUPPORTED, "")
}

func (app *WebApp) handleEvent(ic iris.Context, txn newrelic.Transaction) {
	request := ic.Request()
	now := CurrentTimestamp()

	rawData, err := ioutil.ReadAll(ic.Request().Body)
	if err != nil && err != io.ErrUnexpectedEOF {
		txn.NoticeError(err)
		txn.AddAttribute("http-response-status-code", 500)
		txn.AddAttribute("errer-stmt", err.Error())
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
		txn.AddAttribute("errer-stmt", err.Error())
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
		txn.AddAttribute("errer-stmt", err.Error())
		WriteError(ic, 500, EXCEPTION_MSG_GENERAL, err.Error())
		return
	}

	txn.AddAttribute("eventCategory", eventCategory)

	decoded := map[string]interface{}{}
	if err := json.Unmarshal(rawData, &decoded); err != nil {
		txn.NoticeError(err)
		txn.AddAttribute("http-response-status-code", 500)
		txn.AddAttribute("errer-stmt", err.Error())
		WriteError(ic, 500, EXCEPTION_MSG_GENERAL, err.Error())
		return
	}

	// assign clientIP
	if dm := reflect.ValueOf(decoded["device"]); dm.Kind() == reflect.Map {
		dm.SetMapIndex(reflect.ValueOf("clientIP"), reflect.ValueOf(clientIP))
	}

	// assign recvTimestamp
	decoded["recvTimestamp"] = now

	mobileEvent := MobileEvent{}
	if err := json.Unmarshal(rawData, &mobileEvent); err != nil {
		txn.NoticeError(err)
		txn.AddAttribute("http-response-status-code", 500)
		txn.AddAttribute("errer-stmt", err.Error())
		WriteError(ic, 500, EXCEPTION_MSG_GENERAL, err.Error())
		return
	}

	payload := EventLog{
		WhatToDo:      JOB_NAME_MOBILE_EVENT, // what_to_do
		LogUUID:       logUUID.String(),      // log_uuid
		RecvTimestamp: now,                   // recv_timestamp
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
		txn.AddAttribute("errer-stmt", err.Error())
		WriteError(ic, 500, EXCEPTION_MSG_GENERAL, err.Error())
		return
	}

	if err := app.producer.Publish("airbridge-raw-events", pk, encoded); err != nil {
		txn.NoticeError(err)
		txn.AddAttribute("http-response-status-code", 500)
		txn.AddAttribute("errer-stmt", err.Error())
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

	log.Printf("[200][%s] app: %s, event_category: %d", clientIP, appName, eventCategory)
}
