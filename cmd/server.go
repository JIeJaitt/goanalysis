package cmd

import (
	"flag"
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/spf13/cobra"
	"github.com/toheart/goanalysis/internal/conf"
	"github.com/toheart/goanalysis/internal/pkg/logger"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "./configs", "config path, eg: -conf config.yaml")
	rootCmd.AddCommand(ServerCmd)
}

func newApp(logger log.Logger, servers []transport.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(servers...),
	)
}

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "start the server",
	Run: func(cmd *cobra.Command, args []string) {
		c := config.New(
			config.WithSource(
				file.NewSource(flagconf),
			),
		)
		defer c.Close()

		if err := c.Load(); err != nil {
			panic(err)
		}

		var bc conf.Bootstrap
		if err := c.Scan(&bc); err != nil {
			panic(err)
		}

		// 启动服务器的逻辑
		log := logger.NewLogger(bc.Logger)
		app, cleanup, err := wireApp(bc.Server, bc.Biz, log)
		if err != nil {
			panic(err)
		}
		defer cleanup()

		// start and wait for stop signal
		if err := app.Run(); err != nil {
			panic(err)
		}
	},
}
