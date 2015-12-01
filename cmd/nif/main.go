package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/scraperwiki/nif"
)

var version = "?"

func init() {
	log.SetFlags(0)
}

func main() {
	app := cli.NewApp()
	app.Name = "nif"
	app.Usage = "Simple network interface information tool"
	app.Version = version
	app.Action = actionMain

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "all, a",
			Usage: "List all available network interfaces",
		},
		cli.BoolFlag{
			Name:  "one, o, 1",
			Usage: "Show only single best guessed network interfaces and/or IP address",
		},
		cli.BoolFlag{
			Name:  "ipv4, 4",
			Usage: "Show IPv4 addresses next to network interface",
		},
		cli.BoolFlag{
			Name:  "ipv6, 6",
			Usage: "Show IPv6 addresses next to network interface",
		},
		cli.BoolFlag{
			Name:  "only-ip, i",
			Usage: "Only show IP addresses of network interface",
		},
		cli.IntFlag{
			Name:  "retry, r",
			Usage: "Retry n times in intervals of 1sec if no interface addresses could be found",
		},
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "Show additional debug information",
		},
	}

	app.RunAndExitOnError()
}

func actionMain(c *cli.Context) {
	if c.Bool("one") && c.Bool("all") {
		log.Fatalln("Error: conflicting flags: -1/-o/--one and -a/--all")
	}
	if c.Bool("only-ip") && !c.Bool("ipv4") && !c.Bool("ipv6") {
		log.Fatalln("Error: missing flag: -4/--ipv4 or -6/--ipv6")
	}

	nifs, err := nif.Interfaces(!c.Bool("all"))
	if err != nil {
		log.Fatalln("Error:", err)
	}

	if c.Bool("one") {
		nifs = nifs[0:1]
	}

	for _, n := range nifs {
		v4IPs, v6IPs, err := nif.IPs(n, c.Bool("ipv4"), c.Bool("ipv6"), c.Int("retry"))
		if err != nil {
			log.Fatalln("Error:", err)
		}

		var buf bytes.Buffer

		if !c.Bool("only-ip") {
			if !c.Bool("debug") {
				fmt.Fprintf(&buf, "%s ", n.Name)
			} else {
				fmt.Fprintf(&buf, "%d %s %s %s ",
					n.Index, n.Name, stringHardwareAddress(n.HardwareAddr), n.Flags)
			}
		}

		if c.Bool("ipv4") {
			if len(v4IPs) > 0 {
				if c.Bool("one") {
					v4IPs = v4IPs[0:1]
				}
				fmt.Fprint(&buf, stringIPs(v4IPs))
			}
		}

		if c.Bool("ipv4") && c.Bool("ipv6") && len(v6IPs) > 0 {
			fmt.Fprint(&buf, " ")
		}

		if c.Bool("ipv6") {
			if len(v6IPs) > 0 {
				if c.Bool("one") {
					v6IPs = v6IPs[0:1]
				}
				fmt.Fprint(&buf, stringIPs(v6IPs))
			}
		}

		if len(buf.Bytes()) == 0 {
			continue
		}

		fmt.Fprint(&buf, "\n")

		_, _ = buf.WriteTo(os.Stdout)
	}
}

func stringIPs(addrs []net.IP) string {
	var s string
	for _, addr := range addrs {
		s += addr.String() + "|"
	}

	return strings.TrimRight(s, "|")
}

func stringHardwareAddress(hardwareAddr net.HardwareAddr) string {
	s := hardwareAddr.String()
	return s
}
