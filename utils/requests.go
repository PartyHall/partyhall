package utils

import (
	"net"

	"github.com/labstack/echo/v4"
)

func IsRemote(c *echo.Context) bool {
	// addr := r.Header.Get("X-Real-IP")
	addr := (*c).Request().Header.Get("X-Real-IP")
	if len(addr) == 0 {
		ip, _, err := net.SplitHostPort((*c).RealIP())
		if err != nil {
			// Just to be on the safe side, we say that it's a remote connection when we can't parse the hostport (shouldn't happen)
			return true
		}

		addr = ip
	}

	return addr != "127.0.0.1" && addr != "[::1]"
}
