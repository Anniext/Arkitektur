package utils

import (
	"github.com/gin-gonic/gin"
	"net"
)

func GetLocalIP() string {
	addrList, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrList {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String()
		}
	}
	return "未知"
}

func GetClientIP(ctx *gin.Context) string {
	clientIP := ctx.Request.Header.Get("X-Forwarded-For")
	if clientIP == "" {
		clientIP = ctx.Request.Header.Get("X-Real-IP")
	}
	if clientIP == "" {
		clientIP = ctx.ClientIP()
	}

	return clientIP
}
