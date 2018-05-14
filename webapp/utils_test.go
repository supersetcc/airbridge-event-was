package webapp

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeToStr(t *testing.T) {
	from := time.Date(2018, time.May, 8, 17, 53, 9, 0, time.UTC)
	result := TimeToStr(from)
	if result != "2018-05-08 17:53:09" {
		t.Fatalf("could not convert to expected: %s", result)
	}
}

func TestGetSimpleLinkString(t *testing.T) {
	appSubdomain := "example.com"
	channel := "channel"
	params := map[string]string{
		"hello": "world",
	}

	expected := fmt.Sprintf("http://abr.ge/@%s/%s?hello=world", appSubdomain, channel)
	result, err := GetSimpleLinkString(appSubdomain, channel, params)
	if err != nil || result != expected {
		t.Fatalf("could not get simplelink")
	}
}

func TestGetSimpleLinkStringWithShortID(t *testing.T) {
	shortID := "1234567890asdfghjkl"
	expected := fmt.Sprintf("http://abr.ge/%s", shortID)
	result, err := GetSimpleLinkStringWithShortID(shortID)
	if err != nil || result != expected {
		t.Fatalf("could not get shortid's simplelink")
	}
}

func TestGenerateKafkaPartitionKey(t *testing.T) {
	osVersion := "android"
	deviceModel := "samsung"
	appSubdomain := "google.com"
	remoteAddr := "8.8.8.8"

	expected := "a921196c07729070f531cada955afc3e"
	result := GenerateKafkaPartitionKey(osVersion, deviceModel, appSubdomain, remoteAddr)
	if result != expected {
		t.Fatalf("could not generate kafka's partition key")
	}
}
