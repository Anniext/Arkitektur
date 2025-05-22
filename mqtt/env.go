package mqtt

import "os"

var idc string
var podIp string
var nodeIp string

func GetIDC() string {
	idcEnv := os.Getenv("IDC")
	if len(idcEnv) == 0 {
		return "unknown"
	} else {
		idc = idcEnv
	}

	return idc
}

func GetPodIp() string {
	podEnv := os.Getenv("POD_IP")
	if len(podEnv) == 0 {
		return "unknown"
	} else {
		podIp = podEnv
	}

	return podIp
}

func GetNodeIp() string {
	nodeEnv := os.Getenv("NODE_IP")
	if len(nodeEnv) == 0 {
		return "unknown"
	} else {
		nodeIp = nodeEnv
	}

	return nodeIp
}
