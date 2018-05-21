package webapp

import (
	"testing"

	"bitbucket.org/teamteheranslippers/airbridge-go-bypass-was/common"
	"github.com/iris-contrib/httpexpect"
	"github.com/kataras/iris"
	"github.com/kataras/iris/httptest"
)

func NewWebAppExpect(t *testing.T) (*httpexpect.Expect, *common.MessageProducerMock) {
	app := iris.New()
	mp := &common.MessageProducerMock{false, "", "", nil}
	logging, err := common.NewLoggingDebug()
	if err != nil {
		t.Fatalf("could not get LoggingDebug: %v", err)
	}

	_, err = NewWebApp(app, mp, logging)
	if err != nil {
		t.Fatalf("could not allocate a webapp: %v", err)
	}

	return httptest.New(t, app), mp
}
