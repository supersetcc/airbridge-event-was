package common

import (
	newrelic "github.com/newrelic/go-agent"
)

type Logging interface {
	NewTransaction() (newrelic.Transaction, error)
}
