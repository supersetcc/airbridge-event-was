package common

import newrelic "github.com/newrelic/go-agent"

type LoggingDebug struct {
	newrelic.Application
}

func NewLoggingDebug() (*LoggingDebug, error) {
	config := newrelic.NewConfig("", "")
	config.Enabled = false

	na, err := newrelic.NewApplication(config)
	if err != nil {
		return nil, err
	}

	return &LoggingDebug{na}, nil
}

func (logging *LoggingDebug) NewTransaction() (newrelic.Transaction, error) {
	transaction := logging.StartTransaction("", nil, nil)
	return transaction, nil
}
