package data

import "log"

type DBConfig struct {
	Mode   string
	Dns    string
	Driver string
}

var defaultDB *DBConfig

func GetDefaultDB() *DBConfig {
	return defaultDB
}

func InitDefaultDB(mode, dns, driver string) error {
	defaultDB = &DBConfig{
		Mode:   mode,
		Dns:    dns,
		Driver: driver,
	}

	err := Init()
	if err != nil {
		log.Panic(err.Error())
		return err
	}
	err = Run()
	if err != nil {
		log.Panic(err.Error())
		return err
	}
	return nil
}
