package mqtt

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"runtime/debug"
)

// 生成随机客户端ID
func generateClientID() string {
	podIp := GetPodIp()
	nodeIp := GetNodeIp()
	iDC := GetIDC()

	if len(podIp) != 0 && len(iDC) != 0 && len(nodeIp) != 0 {
		return fmt.Sprintf("%s@%s@%s", iDC, nodeIp, podIp)
	} else {
		b := make([]byte, 4)
		rand.Read(b)
		return fmt.Sprintf("%s", hex.EncodeToString(b))
	}
}

// SafeGoRecoverWarpFunc function    安全运行协程
func SafeGoRecoverWarpFunc(h func()) func() {
	return func() {
		var err error
		defer func() {
			r := recover()
			if r != nil {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("unkonw error")
				}

				log.Println(err.Error())
				log.Println("stack: ", string(debug.Stack()))
			}

		}()

		h()
	}
}
