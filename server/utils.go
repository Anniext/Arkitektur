package server

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
