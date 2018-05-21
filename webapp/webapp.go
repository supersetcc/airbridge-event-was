package webapp

import (
	"bitbucket.org/teamteheranslippers/airbridge-go-bypass-was/common"
	"github.com/kataras/iris"
)

type WebApp struct {
	producer common.MessageProducer
	logging  common.Logging
}

func NewWebApp(app *iris.Application, mp common.MessageProducer, logging common.Logging) (*WebApp, error) {
	webapp := &WebApp{mp, logging}
	// handle mobile event
	app.Post("/api/v2/apps/{app_name}/events/mobile-app/{event_category}", webapp.HandleMobileEventReceiver)

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
