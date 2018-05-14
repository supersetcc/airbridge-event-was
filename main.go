package main

import (
	"context"
	"log"
	"time"

	webapp "bitbucket.org/teamteheranslippers/airbridge-go-stat-udl-io/webapp"
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

	// to support graceful shutdown, iris support to catch a Interrupt
	iris.RegisterOnInterrupt(func() {
		log.Printf("shutdown airbridge-go-stat-udl-io")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
		defer cancel()

		app.Shutdown(ctx)
	})

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

	_, err = webapp.NewWebApp(app, config.Kafka.BrokerList)
	app.Run(iris.Listener(listener), iris.WithoutInterruptHandler)
}
