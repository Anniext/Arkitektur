package binance

import (
	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
)

var defaultBinance *binance.Client
var defaultFutures *futures.Client

func InitBinance() error {
	cnf := GetDefaultBinanceConfig()

	if cnf.Proxy == "" {
		defaultBinance = binance.NewClient(cnf.ApiKey, cnf.SecretKey)
		defaultFutures = futures.NewClient(cnf.ApiKey, cnf.SecretKey)
	} else {
		defaultBinance = binance.NewProxiedClient(cnf.ApiKey, cnf.SecretKey, cnf.Proxy)
		defaultFutures = futures.NewProxiedClient(cnf.ApiKey, cnf.SecretKey, cnf.Proxy)
	}

	return nil
}

func GetDefaultBinance() *binance.Client {
	return defaultBinance
}

func GetDefaultFutures() *futures.Client {
	return defaultFutures
}
