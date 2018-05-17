package webapp

import (
	"bitbucket.org/teamteheranslippers/airbridge-go-bypass-was/common"
	"github.com/kataras/iris"
)

type WebApp struct {
	producer common.MessageProducer
}

func NewWebApp(app *iris.Application, mp common.MessageProducer) (*WebApp, error) {
	webapp := &WebApp{mp}
	app.Post("/api/v2/apps/{app_name}/events/mobile-app/{event_category}", webapp.HandleMobileEventReceiver)

	// handle unsupported method
	app.Delete("/api/v2/apps/{app_name}/events/mobile-app/{event_category}", webapp.HandleUnsupportedMethod)
	app.Get("/api/v2/apps/{app_name}/events/mobile-app/{event_category}", webapp.HandleUnsupportedMethod)
	app.Head("/api/v2/apps/{app_name}/events/mobile-app/{event_category}", webapp.HandleUnsupportedMethod)
	app.Options("/api/v2/apps/{app_name}/events/mobile-app/{event_category}", webapp.HandleUnsupportedMethod)
	app.Patch("/api/v2/apps/{app_name}/events/mobile-app/{event_category}", webapp.HandleUnsupportedMethod)

	return webapp, nil
}

func (webapp *WebApp) Close() error {
	return webapp.producer.Close()
}
