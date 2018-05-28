package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	common "bitbucket.org/teamteheranslippers/airbridge-go-bypass-was/common"
	webapp "bitbucket.org/teamteheranslippers/airbridge-go-bypass-was/webapp"
	cors "github.com/iris-contrib/middleware/cors"
	iris "github.com/kataras/iris"
	tcplisten "github.com/valyala/tcplisten"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// config
	production := len(os.Getenv("USE_AIRBRIDGE_LOCAL_DB")) == 0
	if production == true {
		log.Printf("running on production environment")
	}

	config, err := LoadConfig(production)
	if err != nil {
		log.Fatalf("could not config load: %v", err)
	}

	// iris init
	app := iris.New()

	// package tcplisten provides customizeable TCP net.Listener with various
	// performance-related options
	listenerConfig := tcplisten.Config{
		ReusePort:   true,
		DeferAccept: true,
		FastOpen:    true,
	}

	// allow all origins, allow methods: GET and POST
	app.Use(cors.Default())

	listener, err := listenerConfig.NewListener("tcp4", fmt.Sprintf(":%d", config.Server.Port))
	if err != nil {
		log.Fatalf("could not open socket from tcplisten: %v", err)
	}

	mp, err := common.NewMessageProducerKafka(config.Kafka.BrokerList)
	if err != nil {
		log.Fatalf("could not open kafka producer: %v", err)
	}

	logger, err := common.NewLoggingNewrelic(config.Newrelic.AppName, config.Newrelic.License)
	if err != nil {
		log.Fatalf("could not open LoggingNewrelic: %v", err)
	}

	wa, err := webapp.NewWebApp(app, mp, logger)
	if err != nil {
		log.Fatalf("could not allocate a WebApp: %v", err)
	}

	// to support graceful shutdown, iris support to catch a Interrupt
	iris.RegisterOnInterrupt(func() {
		log.Printf("shutdown airbridge-go-bypass-was")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
		defer cancel()

		app.Shutdown(ctx)

		if err := wa.Close(); err != nil {
			log.Printf("close error: %v", err)
		}
	})

	app.Run(iris.Listener(listener), iris.WithoutInterruptHandler)
}
