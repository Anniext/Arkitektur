package middlewares

import (
	"github.com/Anniext/Arkitektur/casbin"
	"github.com/Anniext/Arkitektur/code"
	"github.com/Anniext/Arkitektur/common"
	"github.com/Anniext/Arkitektur/jwt"
	"github.com/Anniext/Arkitektur/system/log"
	"github.com/Anniext/Arkitektur/utils"
	"github.com/gin-gonic/gin"
	"strings"
)

var CasbinHandler = Casbin()

func Casbin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if strings.ToLower(ctx.Request.Method) == "OPTIONS" {
			ctx.Next()
			return
		}

		api := &common.CodeApi{}
		roleId := jwt.GetTokenData[int64](ctx, "role_id")
		if roleId == 0 {
			log.Errorln("无法使用 `Casbin` 权限校验, 请确保 `Token` 中包含了 `role_id`")
			ctx.Next()
			return
		}

		path := utils.ConvertToRestfulURL(strings.TrimPrefix(ctx.Request.URL.Path, "/api"))

		success, err := casbin.GetDefaultCasbin().Enforce(utils.Int64ToString(roleId), path, ctx.Request.Method)
		if err != nil {
			api.Fail(ctx, code.ErrCodeCasbinNotActiveYet)
			return
		}
		if success {
			ctx.Next()
		} else {
			ctx.Abort()
			api.Fail(ctx, code.ErrCodeCasbinNotPermissions)
			return
		}
	}
}
