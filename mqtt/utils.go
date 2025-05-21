package mqtt

import (
	"admin/common/env"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// 生成随机客户端ID
func generateClientID() string {
	podIp := env.GetPodIp()
	nodeIp := env.GetNodeIp()
	iDC := env.GetIDC()

	if len(podIp) != 0 && len(iDC) != 0 && len(nodeIp) != 0 {
		return fmt.Sprintf("%s@%s@%s", iDC, nodeIp, podIp)
	} else {
		b := make([]byte, 4)
		rand.Read(b)
		return fmt.Sprintf("unknown@%s", hex.EncodeToString(b))
	}
}
