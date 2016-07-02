package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("cfx", "Cloudformation Toolkit")

	showParameters          = app.Command("show-parameters", "Show parameters for a stack.")
	showParametersStackName = showParameters.Flag("stack-name", "Name of the stack.").Required().String()

	updateParameters             = app.Command("update-parameters", "Update parameters for a stack.")
	updateParametersStackName    = updateParameters.Flag("stack-name", "Name of the stack.").Required().String()
	updateParametersCapabilities = updateParameters.Flag("capabilities", "Capabilities required to perform changes.").Strings()
	updateParametersParameters   = updateParameters.Flag("parameter", "Parameters to change.").StringMap()

	settle             = app.Command("settle", "Wait for a stack to settle, tailing events.")
	settleStackName    = settle.Flag("stack-name", "Name of the stack.").Required().String()
	settleTimeout      = settle.Flag("timeout", "Maximum time to wait until the stack is considered settled.").Default("10m").Duration()
	settlePollInterval = settle.Flag("poll-interval", "Interval at which to poll AWS for events").Default("5s").Duration()
)

func main() {
	c := kingpin.MustParse(app.Parse(os.Args[1:]))

	awsSession := session.New()

	cf := cloudformation.New(awsSession)

	switch c {
	case showParameters.FullCommand():
		if err := handleShowParameters(cf, *showParametersStackName); err != nil {
			panic(err)
		}
	case updateParameters.FullCommand():
		if err := handleUpdateParameters(cf, *updateParametersStackName, *updateParametersCapabilities, *updateParametersParameters); err != nil {
			panic(err)
		}
	case settle.FullCommand():
		if err := handleSettle(cf, *settleStackName, *settleTimeout, *settlePollInterval); err != nil {
			panic(err)
		}
	}
}
