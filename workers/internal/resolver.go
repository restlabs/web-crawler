package internal

import (
	"errors"
	"fmt"
	"net"
)

var (
	ErrorIpCannotBeResolved = errors.New("invalid URL no IPs to resolve")
)

func ResolveDns(host string) error {
	ips, err := net.LookupIP(host)
	if err != nil {
		return fmt.Errorf("error resolving IP for host %v. %v", host, err)
	}

	if len(ips) == 0 {
		return ErrorIpCannotBeResolved
	}

	return nil
}
