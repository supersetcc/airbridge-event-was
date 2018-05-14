package webapp

import (
	"github.com/kataras/iris"
)

func WriteError(ic iris.Context, code int, message string, hint string) {
	ic.StatusCode(code)
	ic.Header("Cache-Control", "no-cache")
	ic.Header("Pragma", "no-cache")

	at := TimeToStr(KSTNow())
	response := map[string]interface{}{
		"resultMessage": message,
		"hint":          hint,
		"at":            at,
	}

	ic.JSON(response)
}

func WriteResponse(ic iris.Context, response interface{}) {
	ic.StatusCode(200)
	ic.Header("Cache-Control", "no-cache")
	ic.Header("Pragma", "no-cache")

	ic.JSON(response)
}
