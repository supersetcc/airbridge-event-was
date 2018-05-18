package webapp

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"bitbucket.org/teamteheranslippers/airbridge-go-bypass-was/common"
	"github.com/iris-contrib/httpexpect"
	"github.com/kataras/iris"
	"github.com/kataras/iris/httptest"
)

type MockMessageProducer struct {
	IsClosed                  bool
	LastPublishedTopic        string
	LastPublishedPartitionKey string
	LastPublishedPayload      []byte
}

func (p *MockMessageProducer) Publish(topic, pk string, payload []byte) error {
	p.LastPublishedTopic = topic
	p.LastPublishedPartitionKey = pk
	p.LastPublishedPayload = payload
	return nil
}

func (p *MockMessageProducer) Close() error {
	p.IsClosed = true
	return nil
}

func MakeWebAppExpect(t *testing.T) (*httpexpect.Expect, *MockMessageProducer) {
	app := iris.New()
	mp := &MockMessageProducer{false, "", "", nil}
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

func TestMockMobileEventRequestWithoutAhtorized(t *testing.T) {
	expect, _ := MakeWebAppExpect(t)
	uri := fmt.Sprintf("/api/v2/apps/%s/events/mobile-app/%d", "ablog", 9162)
	expect.POST(uri).Expect().Status(httptest.StatusUnauthorized)
}

func TestMockMobileEventRequestWithAuthorized(t *testing.T) {
	expect, _ := MakeWebAppExpect(t)
	uri := fmt.Sprintf("/api/v2/apps/%s/events/mobile-app/%d", "ablog", 9162)
	expect.POST(uri).WithHeader("Authorization", "random-authorized-string").Expect().Status(httptest.StatusBadRequest)
}

func TestMockMobileEventRequestBasic(t *testing.T) {
	path := "../res/test/mobile_event_request_data.txt"
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		t.Fatalf("could not open test data(%s): %v", path, err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		expect, mp := MakeWebAppExpect(t)
		uri := fmt.Sprintf("/api/v2/apps/%s/events/mobile-app/%d", "ablog", 9162)

		request := expect.POST(uri).WithHeader("Authorization", "random-authorized-string")
		payload := scanner.Text()
		request.WithText(payload)
		request.Expect().Status(httptest.StatusOK)

		if mp.LastPublishedTopic != "airbridge-raw-events" {
			t.Fatalf("publish is not sent to 'airbrdige-raw-events' but %v", mp.LastPublishedTopic)
		}
	}
}
