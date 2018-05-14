package webapp

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"
)

func KSTNow() time.Time {
	utc := time.Now().UTC()
	return utc.Add(time.Hour * 9)
}

func TimeToStr(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

func CurrentTimestamp() int64 {
	return time.Now().UnixNano() / 1000000
}

func GetSimpleLinkString(appSubdomain string, channel string, params map[string]string) (string, error) {
	u, err := url.Parse("http://abr.ge/")
	if err != nil {
		return "", err
	}

	u.Path = fmt.Sprintf("@%s/%s", appSubdomain, channel)

	v := url.Values{}
	for key, value := range params {
		if value == "" {
			continue
		}

		v.Add(key, value)
	}

	u.RawQuery = v.Encode()

	return u.String(), nil
}

func GetSimpleLinkStringWithShortID(shortID string) (string, error) {
	return fmt.Sprintf("http://abr.ge/%s", shortID), nil
}

func GenerateKafkaPartitionKey(osVersion, deviceModel, appSubdomain, remoteAddr string) string {
	hash := md5.New()
	io.WriteString(hash, fmt.Sprintf("%s-%s-%s-%s", remoteAddr, osVersion, deviceModel, appSubdomain))
	return hex.EncodeToString(hash.Sum(nil))
}

func ParseClientIPFromXForwarededFor(forwardedIP string) string {
	if forwardedIP == "" {
		return ""
	}

	segments := strings.Split(forwardedIP, ",")
	return segments[len(segments)-1]
}

func GetCanonicalDeviceUUIDAndGenType() (string, string) {
	return "", ""
}
