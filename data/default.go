package data

import "log"

func InitDefaultDB() error {
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
