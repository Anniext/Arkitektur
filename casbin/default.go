package casbin

import (
	"github.com/Anniext/Arkitektur/data"
	"github.com/casbin/casbin/v2"
	xormadapter "github.com/casbin/xorm-adapter/v3"
)

var defaultCasbin *casbin.SyncedEnforcer

// InitCasbin TODO 依赖数据模块， 目前只支持xorm
// InitCasbin 初始化casbin
func InitCasbin() error {
	cnf := GetDefaultCasbinConfig()

	var err error
	modePath := cnf.ModelPath
	a, _ := xormadapter.NewAdapterByEngine(data.GetDB())
	defaultCasbin, err = casbin.NewSyncedEnforcer(modePath, a)
	if err != nil {
		return err
	}

	err = defaultCasbin.LoadPolicy()
	if err != nil {
		return err
	}

	return nil
}

func GetDefaultCasbin() *casbin.SyncedEnforcer {
	return defaultCasbin
}
