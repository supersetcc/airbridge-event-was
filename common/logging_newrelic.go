package common

import newrelic "github.com/newrelic/go-agent"

type LoggingNewrelic struct {
	newrelic.Application
}

func NewLoggingNewrelic(appName, license string) (*LoggingNewrelic, error) {
	config := newrelic.NewConfig(appName, license)
	app, err := newrelic.NewApplication(config)
	if err != nil {
		return nil, err
	}

	return &LoggingNewrelic{app}, nil
}

func (logging *LoggingNewrelic) NewTransaction() (newrelic.Transaction, error) {
	transaction := logging.StartTransaction("airbrdige-go-bypass", nil, nil)
	return transaction, nil
}
