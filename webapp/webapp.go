package webapp

import "github.com/kataras/iris"

type WebApp struct {
	producer *MessageProducer
}

func NewWebApp(app *iris.Application, brokers []string) (*WebApp, error) {
	producer, err := NewMessageProducer(brokers)
	if err != nil {
		return nil, err
	}

	webapp := &WebApp{producer}
	app.Post("/api/v2/apps/{app_name}/events/mobile-app/{event_category}", webapp.HandleMobileEventReceiver)

	// handle unsupported method
	app.Delete("/api/v2/apps/{app_name}/events/mobile-app/{event_category}", webapp.HandleUnsupportedMethod)
	app.Get("/api/v2/apps/{app_name}/events/mobile-app/{event_category}", webapp.HandleUnsupportedMethod)
	app.Head("/api/v2/apps/{app_name}/events/mobile-app/{event_category}", webapp.HandleUnsupportedMethod)
	app.Options("/api/v2/apps/{app_name}/events/mobile-app/{event_category}", webapp.HandleUnsupportedMethod)
	app.Patch("/api/v2/apps/{app_name}/events/mobile-app/{event_category}", webapp.HandleUnsupportedMethod)

	return webapp, nil
}
