package webapp

import (
	"fmt"
	"strings"

	"github.com/ua-parser/uap-go/uaparser"
)

type OSType struct {
	Family  string
	Version string
}

type UserAgent struct {
	OS     OSType
	Device struct {
		Family string
		Brand  string
		Model  string
	}
	Browser struct {
		Family  string
		Version string
	}
}

var (
	parser *uaparser.Parser
)

func convertDeviceModel(brand, model string) string {
	if brand == "LG" && strings.Index(model, "-") == -1 {
		return fmt.Sprintf("LG-%s", model)
	}

	if brand == "Samsung" && strings.Index(model, "API") >= 0 {
		return fmt.Sprintf("%s %s", brand, model)
	}

	return model
}

func NewUserAgent(stmt string) (*UserAgent, error) {

	if parser == nil {
		var err error

		parser, err = uaparser.New("../res/ua_parser/regexes.yaml")
		if err != nil {
			return nil, err
		}
	}

	ret := parser.Parse(stmt)
	ua := UserAgent{}

	ua.OS.Family = ret.Os.Family
	ua.OS.Version = fmt.Sprintf("%s.%s.%s", ret.Os.Major, ret.Os.Minor, ret.Os.Patch)

	ua.Browser.Family = ret.UserAgent.Family
	ua.Browser.Version = fmt.Sprintf("%s.%s.%s", ret.UserAgent.Major, ret.UserAgent.Minor, ret.UserAgent.Patch)

	ua.Device.Family = ret.Device.Family
	ua.Device.Brand = ret.Device.Brand
	ua.Device.Model = convertDeviceModel(ret.Device.Brand, ret.Device.Model)

	return &ua, nil
}

func (os *OSType) String() string {
	return os.Family + os.Version
}
