package utils

import (
	"errors"
	"log"
	"runtime/debug"
)

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

func SafeGoRecoverWarpFuncInt64(h func() int64) func() int64 {
	return func() int64 {
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
					err = errors.New("unknown error")
				}
				log.Println(err.Error())
				log.Println("stack:", string(debug.Stack()))
			}
		}()
		return h()
	}
}
