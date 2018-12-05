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
	PartitionKey  string         `json:"partition_key"`
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

func NewWebApp(app *iris.Application, mp common.MessageProducer, logging common.Logging) (*WebApp, error) {
	webapp := &WebApp{mp, logging}

	// handle mobile event Receiver
	app.Post("/api/v2/apps/{app_name}/events/mobile-app/{event_category}", webapp.HandleEventReceiverMobile)
	app.Post("/api/v3/apps/{app_name}/events/mobile-app/{event_category}", webapp.HandleEventReceiverMobile)
	app.Post("/api/v3.1/apps/{app_name}/events/mobile-app/{event_category}", webapp.HandleEventReceiverMobile)

	// handle web event Receiver
	app.Post("/api/v2/apps/{app_name}/events/mobile-webapp/{event_category}", webapp.HandleEventReceiverWebApp)
	app.Post("/api/v3/apps/{app_name}/events/mobile-webapp/{event_category}", webapp.HandleEventReceiverWebApp)
	app.Post("/api/v3.1/apps/{app_name}/events/mobile-webapp/{event_category}", webapp.HandleEventReceiverWebApp)

	// handle 404 error
	app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		WriteError(ctx, 404, "invalid request", "")
	})

	// handle health check
	app.Get("/health-check", webapp.HandleHealthCheck)

	return webapp, nil
}

func (webapp *WebApp) Close() error {
	return webapp.producer.Close()
}
