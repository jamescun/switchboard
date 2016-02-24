package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/jamescun/switchboard/backend"
	"github.com/jamescun/switchboard/match"
	"github.com/jamescun/switchboard/server"

	"github.com/codegangsta/cli"
)

var (
	BuildVersion  = "development"
	BuildRevision = "development"
)

var App = &cli.App{
	Name:        "switchboard",
	HelpName:    "switchboard",
	Usage:       "tls-based service discovery proxy",
	HideVersion: true,
	Action:      cli.ShowAppHelp,
	Writer:      os.Stdout,

	Commands: []cli.Command{
		{
			Name:  "server",
			Usage: "launch proxy server",

			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "listen",
					Value:  "127.0.0.1:1993",
					Usage:  "host:port for proxy to listen for client requests",
					EnvVar: "SWITCHBOARD_LISTEN",
				},
				cli.DurationFlag{
					Name:   "cache",
					Value:  1 * time.Minute,
					Usage:  "timeout for caching service discovery results",
					EnvVar: "SWITCHBOARD_CACHE",
				},
				cli.StringFlag{
					Name:   "balance",
					Value:  "random",
					Usage:  "name of load balancing scheme (random|consistent|rendezvous)",
					EnvVar: "SWITCHBOARD_BALANCE",
				},
				cli.IntFlag{
					Name:   "buffer-size",
					Value:  4096,
					Usage:  "buffer size for rx/tx between client and upstream",
					EnvVar: "SWITCHBOARD_BUFFER_SIZE",
				},
			},

			Subcommands: []cli.Command{
				{
					Name:  "static",
					Usage: "proxy all to a static set of upstream servers",

					Flags: []cli.Flag{
						cli.StringSliceFlag{
							Name:   "upstream",
							Usage:  "host:port pairs of upstream server addresses",
							EnvVar: "SWITCHBOARD_STATIC_UPSTREAM",
						},
					},

					Action: ServerStatic,
				},
			},
		},
		{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "show version information",
			Action:  Version,
		},
	},
}

func ServerStatic(c *cli.Context) {
	bk := backend.Static{
		Upstreams: c.StringSlice("upstream"),
	}

	Server(c, bk)
}

func Server(c *cli.Context, bk backend.Backend) {
	s := &server.Server{
		Match:      match.Http,
		Backend:    bk,
		BufferSize: c.GlobalInt("buffer-size"),
	}

	ln, err := net.Listen("tcp", c.GlobalString("listen"))
	if err != nil {
		log.Fatalln("fatal: listen:", err)
		return
	}
	log.Println("info: listen:", c.GlobalString("listen"))

	err = s.Serve(ln)
	if err != nil {
		log.Fatalln("fatal: serve:", err)
		return
	}
}

func Version(c *cli.Context) {
	fmt.Fprintln(c.App.Writer, "Build Version: ", BuildVersion)
	fmt.Fprintln(c.App.Writer, "Build Revision:", BuildRevision)
}

func main() {
	App.Run(os.Args)
}
