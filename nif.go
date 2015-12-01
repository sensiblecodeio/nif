package nif

import (
	"net"
	"strings"
	"time"
)

func Interfaces(filtered bool) ([]net.Interface, error) {
	nifs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	if !filtered {
		return nifs, nil
	}

	var nifsFiltered []net.Interface
	for _, nif := range nifs {
		if isOfInterest(nif) {
			nifsFiltered = append(nifsFiltered, nif)
		}
	}

	return nifsFiltered, nil
}

func Addresses(nif net.Interface, v4 bool, v6 bool, retries int) ([]net.Addr, []net.Addr, error) {
	var v4Addrs, v6Addrs []net.Addr

	for n := retries; n >= 0; n-- {
		addrs, err := nif.Addrs()
		if err != nil {
			return nil, nil, err
		}

		v4Addrs, v6Addrs = partition(addrs)

		if v4 && v6 && len(v4Addrs) > 0 && len(v6Addrs) > 0 {
			break
		} else if v4 && len(v4Addrs) > 0 {
			v6Addrs = []net.Addr{}
			break
		} else if v6 && len(v6Addrs) > 0 {
			v4Addrs = []net.Addr{}
			break
		}

		if n > 0 {
			time.Sleep(1 * time.Second)
		}
	}

	return v4Addrs, v6Addrs, nil
}

func IPs(nif net.Interface, v4 bool, v6 bool, retries int) ([]net.IP, []net.IP, error) {
	v4Addrs, v6Addrs, err := Addresses(nif, v4, v6, retries)
	if err != nil {
		return nil, nil, err
	}

	v4IPs, err := extractIPs(v4Addrs)
	if err != nil {
		return nil, nil, err
	}

	v6IPs, err := extractIPs(v6Addrs)
	if err != nil {
		return nil, nil, err
	}

	return v4IPs, v6IPs, nil
}

func isOfInterest(nif net.Interface) bool {
	// Filter out interfaces that have no hardware address.
	if len(nif.HardwareAddr) == 0 {
		return false
	}
	// Filter out loopback interfaces.
	if nif.Flags&net.FlagLoopback == net.FlagLoopback {
		return false
	}
	// Filter out Point-to-Point interfaces.
	if nif.Flags&net.FlagPointToPoint == net.FlagPointToPoint {
		return false
	}
	// Filter out interfaces that are not up.
	if nif.Flags&net.FlagUp != net.FlagUp {
		return false
	}

	// If the interface made it that far, assume it's of interest.
	return true
}

func partition(addrs []net.Addr) ([]net.Addr, []net.Addr) {
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

func isIPv4(addr net.Addr) bool {
	return strings.Contains(addr.String(), ".")
}

func extractIPs(addrs []net.Addr) ([]net.IP, error) {
	var IPs []net.IP

	for _, addr := range addrs {
		IP, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			return nil, err
		}
		IPs = append(IPs, IP)
	}

	return IPs, nil
}
