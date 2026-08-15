package main

import (
	"fmt"
	"go_web_template/routes"
	"os"

	mingfuflags "github.com/Fu-XDU/mingfu_go_common/flags"
	"github.com/labstack/gommon/log"
	"github.com/urfave/cli/v2"
	// MYSQL支持：1/3
	// "go_web_template/database"
)

const (
	clientIdentifier = "go_web_template"
	clientUsage      = "go_web_template"
)

var (
	// Injected by -ldflags at build time.
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"

	app = cli.NewApp()
)

func formatAppVersion() string {
	return fmt.Sprintf("%s\n   commit:     %s\n   build time: %s", Version, Commit, BuildTime)
}

func init() {
	app.Action = ServerApp
	app.Name = clientIdentifier
	app.Version = formatAppVersion()
	app.Usage = clientUsage
	app.Commands = []*cli.Command{}
	app.Flags = append(app.Flags, mingfuflags.GinFlags...)
	// MYSQL支持：2/3
	// app.Flags = append(app.Flags, mingfuflags.MysqlFlags...)

	cli.VersionPrinter = func(c *cli.Context) {
		_, _ = fmt.Fprintf(c.App.Writer, "%s\n  version:    %s\n  commit:     %s\n  build time: %s\n",
			c.App.Name, Version, Commit, BuildTime)
	}
}

func ServerApp(ctx *cli.Context) error {
	if args := ctx.Args(); args.Len() > 0 {
		return fmt.Errorf("invalid command: %q", args.First())
	}
	log.Infof("starting %s version=%s commit=%s build_time=%s", clientIdentifier, Version, Commit, BuildTime)
	err := prepare()
	if err != nil {
		log.Error(err)
	}
	return err
}

func prepare() (err error) {
	// MYSQL支持：3/3
	// if err = database.Setup(); err != nil {
	// 	return
	// }
	routes.Run()
	return
}

func main() {
	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
