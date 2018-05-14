package webapp

import (
	"testing"
)

const (
	UA_IPHONE  = "Mozilla/5.0 (iPhone; CPU iPhone OS 6_1_4 like Mac OS X) AppleWebKit/536.26 (KHTML, like Gecko) Version/6.0 Mobile/10B350 Safari/8536.25"
	UA_IPAD    = "Mozilla/5.0 (iPad; CPU OS 8_1_1 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) CriOS/39.0.2171.50 Mobile/12B435 Safari/600.1.4 (000835)"
	UA_ANDROID = "Mozilla/5.0 (Linux; U; Android 2.1-update1; ko-kr; Nexus One Build/ERE27) AppleWebKit/530.17 (KHTML, like Gecko) Version/4.0 Mobile Safari/530.17"
)

func TestUserAgentForIphone(t *testing.T) {
	ua, err := NewUserAgent(UA_IPHONE)
	if err != nil {
		t.Fatalf("could not parse UA_IPHONE: %v", err)
	}

	if ua.OS.String() != "iOS6.1.4" {
		t.Fatalf("invalid OS.String(): %s", ua.OS.String())
	}

	if ua.Device.Model != "iPhone" {
		t.Fatalf("invalid device model: %s", ua.Device.Model)
	}
}

func TestUserAgentIpad(t *testing.T) {
	ua, err := NewUserAgent(UA_IPAD)
	if err != nil {
		t.Fatalf("could not parse UA_IPHONE: %v", err)
	}

	if ua.OS.String() != "iOS8.1.1" {
		t.Fatalf("invalid OS.String(): %s", ua.OS.String())
	}

	if ua.Device.Model != "iPad" {
		t.Fatalf("invalid device model: %s", ua.Device.Model)
	}
}

func TestUserAgentAndroid(t *testing.T) {
	ua, err := NewUserAgent(UA_ANDROID)
	if err != nil {
		t.Fatalf("could not parse UA_IPHONE: %v", err)
	}

	if ua.OS.String() != "Android2.1.update1" {
		t.Fatalf("invalid OS.String(): %s", ua.OS.String())
	}

	if ua.Device.Model != "Nexus One" {
		t.Fatalf("invalid device model: %s", ua.Device.Model)
	}
}
