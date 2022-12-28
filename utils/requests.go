package utils

import (
	"fmt"
	"net"
	"net/http"
)

func IsRemote(r *http.Request) bool {
	addr := r.Header.Get("X-Real-IP")
	if len(addr) == 0 {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			fmt.Println(err)

			// Just to be on the safe side, we say that it's a remote connection when we can't parse the hostport (shouldn't happen)
			return true
		}

		addr = ip
	}

	return addr != "127.0.0.1" && addr != "[::1]"
}
