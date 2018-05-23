package webapp

import (
	"bitbucket.org/teamteheranslippers/airbridge-go-bypass-was/common"
	"github.com/kataras/iris"
)

type WebApp struct {
	producer common.MessageProducer
	logging  common.Logging
}

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
	RecvTimestamp int64          `json:"recv_timestamp"`
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

func NewWebApp(app *iris.Application, mp common.MessageProducer, logging common.Logging) (*WebApp, error) {
	webapp := &WebApp{mp, logging}

	// handle mobile event
	app.Post("/api/v2/apps/{app_name}/events/mobile-app/{event_category}", webapp.HandleEventReceiverMobile)

	// handle unsupported method
	app.Delete("/api/v2/apps/{app_name}/events/mobile-app/{event_category}", webapp.HandleUnsupportedMethod)
	app.Get("/api/v2/apps/{app_name}/events/mobile-app/{event_category}", webapp.HandleUnsupportedMethod)
	app.Head("/api/v2/apps/{app_name}/events/mobile-app/{event_category}", webapp.HandleUnsupportedMethod)
	app.Options("/api/v2/apps/{app_name}/events/mobile-app/{event_category}", webapp.HandleUnsupportedMethod)
	app.Patch("/api/v2/apps/{app_name}/events/mobile-app/{event_category}", webapp.HandleUnsupportedMethod)

	// handle health check
	app.Get("/health-check", webapp.HandleHealthCheck)

	return webapp, nil
}

func (webapp *WebApp) Close() error {
	return webapp.producer.Close()
}
