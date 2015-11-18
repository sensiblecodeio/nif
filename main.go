package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/codegangsta/cli"
)

var version = "?"

func init() {
	log.SetFlags(0)
}

func main() {
	app := cli.NewApp()
	app.Name = "nif"
	app.Usage = "Simple network interface info tool"
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

	nifs, err := net.Interfaces()
	if err != nil {
		log.Fatalln("Error:", err)
	}

	for _, nif := range nifs {
		if !c.Bool("all") &&
			(len(nif.HardwareAddr) == 0 ||
				nif.Flags&net.FlagLoopback == net.FlagLoopback ||
				nif.Flags&net.FlagPointToPoint == net.FlagPointToPoint ||
				nif.Flags&net.FlagUp != net.FlagUp) {
			continue
		}

		var v4Addrs, v6Addrs []net.Addr

		if c.Bool("ipv4") || c.Bool("ipv6") || c.Bool("only-ip") {
			addrs, err := nif.Addrs()
			if err != nil {
				log.Fatalln("Error:", err)
			}

			v4Addrs, v6Addrs = splitIPs(addrs)
		}

		if !c.Bool("only-ip") {
			if !c.Bool("debug") {
				fmt.Print(nif.Name)
			} else {
				//ips := stringIPs(append(v4Addrs, v6Addrs...))
				fmt.Print(nif.Index, " ", nif.Name, " ", nif.HardwareAddr.String(), " ", nif.Flags.String())
			}
			fmt.Print(" ")
		}

		if c.Bool("ipv4") && len(v4Addrs) > 0 {
			ips := parseIPs(v4Addrs)
			if c.Bool("one") {
				fmt.Print(ips[0])
			} else {
				fmt.Print(stringIPs(ips))
			}
		}

		if c.Bool("ipv6") && len(v6Addrs) > 0 {
			fmt.Print(" ")
			ips := parseIPs(v6Addrs)
			if c.Bool("one") {
				fmt.Print(ips[0])
			} else {
				fmt.Print(stringIPs(ips))
			}
		}

		fmt.Println()

		if c.Bool("one") {
			break
		}
	}
}

func isIPv4(addr net.Addr) bool {
	return strings.Contains(addr.String(), ".")
}

func splitIPs(addrs []net.Addr) ([]net.Addr, []net.Addr) {
	var v4Addrs []net.Addr
	var v6Addrs []net.Addr

	for _, addr := range addrs {
		if isIPv4(addr) {
			v4Addrs = append(v4Addrs, addr)
		} else {
			v6Addrs = append(v6Addrs, addr)
		}
	}

	return v4Addrs, v6Addrs
}

func parseIPs(addrs []net.Addr) []net.IP {
	var ips []net.IP

	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			log.Fatalln("Error:", err)
		}
		ips = append(ips, ip)
	}

	return ips
}

func stringIPs(addrs []net.IP) string {
	var s string
	for _, addr := range addrs {
		s += addr.String() + "|"
	}

	return strings.TrimRight(s, "|")
}