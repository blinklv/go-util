// net.go
//
// Author: blinklv <blinklv@icloud.com>
// Create Time: 2021-05-17
// Maintainer: blinklv <blinklv@icloud.com>
// Last Change: 2021-08-18

package util

import "net"

// GetLocalIP returns local ip address.
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
