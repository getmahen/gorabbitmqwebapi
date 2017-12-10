package actions

import (
	"rabbitmqwebapp/processor"

	"github.com/gobuffalo/buffalo"
)

func PortRequestHandler(processor processor.IProcessor) func(c buffalo.Context) error {
	return func(ctx buffalo.Context) error {
		phoneNumber := ctx.Param("phoneNumber")

		//simulate an outgoing call here
		response, _ := processor.Process(phoneNumber)
		// if err != nil {
		// 	panic(err)
		// }

		return ctx.Render(200, r.JSON(response))
		//return c.Render(200, r.JSON(map[string]string{"message": "Welcome to PortRequestHandler!" + phoneNumber}))
	}
}
