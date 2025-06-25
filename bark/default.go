package bark

import (
	"github.com/Anniext/Arkitektur/system/log"
	"github.com/jzksnsjswkw/go-bark"
)

var defaultBark *bark.Client

func InitDefaultRedis() error {
	cnf := GetDefaultBarkConfig()
	defaultBark = &bark.Client{
		ServerURL: cnf.Host,
	}

	log.Infoln("redis server is running")
	return nil
}

func GetDefaultBark() *bark.Client {
	return defaultBark
}
