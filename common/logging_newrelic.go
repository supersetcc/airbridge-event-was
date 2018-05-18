package common

import newrelic "github.com/newrelic/go-agent"

type LoggingNewrelic struct {
	newrelic.Application

	AppName string
}

func NewLoggingNewrelic(appName, license string) (*LoggingNewrelic, error) {
	config := newrelic.NewConfig(appName, license)
	app, err := newrelic.NewApplication(config)
	if err != nil {
		return nil, err
	}

	return &LoggingNewrelic{app, appName}, nil
}

func (logging *LoggingNewrelic) NewTransaction() (newrelic.Transaction, error) {
	transaction := logging.StartTransaction(logging.AppName, nil, nil)
	return transaction, nil
}
