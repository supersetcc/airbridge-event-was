package main

import (
	"context"
	"log"
	"time"

	common "bitbucket.org/teamteheranslippers/airbridge-go-bypass-was/common"
	webapp "bitbucket.org/teamteheranslippers/airbridge-go-bypass-was/webapp"
	iris "github.com/kataras/iris"
	tcplisten "github.com/valyala/tcplisten"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// config
	config, err := LoadConfig()
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

	listener, err := listenerConfig.NewListener("tcp4", ":8080")
	if err != nil {
		log.Fatalf("could not open socket from tcplisten: %v", err)
	}

	mp, err := common.NewKafkaMessageProducer(config.Kafka.BrokerList)
	if err != nil {
		log.Fatalf("could not open kafka producer: %v", err)
	}

	wa, err := webapp.NewWebApp(app, mp)
	if err != nil {
		log.Fatalf("could not allocate a WebApp: %v", err)
	}

	// to support graceful shutdown, iris support to catch a Interrupt
	iris.RegisterOnInterrupt(func() {
		log.Printf("shutdown airbridge-go-bypass-was")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
		defer cancel()

		if err := wa.Close(); err != nil {
			log.Printf("close error: %v", err)
		}

		app.Shutdown(ctx)
	})

	app.Run(iris.Listener(listener), iris.WithoutInterruptHandler)
}
