package data

import (
	"github.com/Anniext/Arkitektur/system/log"
	"github.com/Anniext/Verktyg/persist/core"
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
	"xorm.io/xorm"
)

// TODO : DB init 暂时只支持mysql

/*
 * 数据库
 * @Description:
 */

var gEngine *xorm.Engine
var engineOnce sync.Once

func GetDB() *xorm.Engine {
	engineOnce.Do(func() {
		var err error
		dbCnf := GetDefaultDB()

		gEngine, err = xorm.NewEngine(dbCnf.Driver, dbCnf.Dns)
		if err != nil {
			gEngine = nil
			log.Info("GetDB error", err)
		} else {
			//gEngine.ShowSQL(true)
			gEngine.SetMaxIdleConns(2)            //设置连接池中的保持连接的最大连接数
			gEngine.SetMaxOpenConns(4)            //设置连接池的打开的最大连接数
			gEngine.SetConnMaxLifetime(time.Hour) //设置连接超时时间
			//gEngine.AddHook()
		}
	})
	return gEngine
}

func Register(name string, persist core.IPersist) {
	core.RegisterPersist(name, persist)
}

func Exit() (err error) {
	log.Infoln("Save Begin")
	core.ExitPersist()
	log.Infoln("Save End")
	log.Infoln("SyncData Begin")
	err = core.SyncDataPersist(true)
	if err != nil {
		return
	}
	log.Infoln("SyncData End")
	return
}

func Run() (err error) {
	log.Infoln("Run Begin")
	err = core.RunPersist()
	if err != nil {
		return
	}
	log.Infoln("Run End")
	return nil
}

func Dead() bool {
	return core.DeadPersist()
}

func Init() (err error) {
	engine := GetDB()
	if engine == nil {
		panic("GetDB Error")
	}
	err = core.SyncPersist()
	if err != nil {
		return
	}
	return
}
