package webapp

import "github.com/kataras/iris"

func (wa *WebApp) HandleHealthCheck(ic iris.Context) {
	ic.StatusCode(200)
	ic.WriteString("ok")
}
