package webapp

import (
	"testing"

	"github.com/kataras/iris/httptest"
)

func TestHealthCheck(t *testing.T) {
	expect, _ := NewWebAppExpect(t)
	expect.GET("/health-check").Expect().Status(httptest.StatusOK)
}
