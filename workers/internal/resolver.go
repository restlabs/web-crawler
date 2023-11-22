package internal

import (
	"fmt"
	"net"
)

func ResolveDns(host string) error {
	ips, err := net.LookupIP(host)
	if err != nil {
		return fmt.Errorf("error resolving IP for host %v. %v", host, err)
	}

	if len(ips) == 0 {
		return fmt.Errorf("invalid host, no IPs found for %v", host)
	}

	return nil
}
