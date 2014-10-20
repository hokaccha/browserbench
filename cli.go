package main

import (
	"os"

	"github.com/codegangsta/cli"
)

const helpTemplate = `
Usage: {{.Name}} [options] url

Options:
{{range .Flags}}  {{.}}
{{end}}`

func NewCliApp() *cli.App {
	cli.AppHelpTemplate = helpTemplate[1:]

	app := cli.NewApp()
	app.Name = "browserbench"
	app.HideHelp = true
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		numberFlag,
		startFlag,
		endFlag,
		browserFlag,
		remoteFlag,
		cli.HelpFlag,
	}
	app.Action = doMain

	return app
}

func doMain(c *cli.Context) {
	args := c.Args()

	if len(args) == 0 {
		cli.ShowAppHelp(c)
		os.Exit(1)
	}

	Benchmark(&BenchmarkOpts{
		Url: args[0],
		Start: c.String("start"),
		End: c.String("end"),
		Number: c.Int("number"),
		Browser: c.String("browser"),
		Remote: c.String("remote"),
	})
}

var numberFlag = cli.IntFlag{
	Name: "number, n",
	Value: 1,
	Usage: "Number of requests",
}

var startFlag = cli.StringFlag{
	Name: "start, s",
	Value: "requestStart",
	Usage: "startMark",
}

var endFlag = cli.StringFlag{
	Name: "end, e",
	Value: "domInteractive",
	Usage: "endMark",
}

var browserFlag = cli.StringFlag{
	Name:  "browser, b",
	Value: "chrome",
	Usage: "browser name",
}

var remoteFlag = cli.StringFlag{
	Name:  "remote, r",
	Value: "http://localhost:4444/wd/hub",
	Usage: "RemoteWebDriver server URL",
}
