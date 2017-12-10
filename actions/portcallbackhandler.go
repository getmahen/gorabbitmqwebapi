package actions

import (
	"rabbitmqwebapp/processor"

	"github.com/gobuffalo/buffalo"
)

func PortCallbackHandler(processor processor.IProcessor) func(c buffalo.Context) error {
	return func(ctx buffalo.Context) error {
		requestId := ctx.Param("requestId")
		if len(requestId) == 0 {
			return ctx.Render(400, nil)
		}

		isResponseSent := processor.SendPortResponse(requestId)

		return ctx.Render(200, r.JSON(isResponseSent))
	}
}
