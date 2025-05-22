package common

import (
	"fmt"
	"github.com/Anniext/Arkitektur/code"
	"github.com/Anniext/Arkitektur/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type CodeApi struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	NowTime int64  `json:"nowTime"`
	UseTime string `json:"useTime"`
}

func (api *CodeApi) Sucess(ctx *gin.Context, data any) {
	c := code.Success
	api.Result(ctx, c.Int32(), data, c.String())
}

func (api *CodeApi) Fail(ctx *gin.Context, errCode code.IErrCode) {
	if errCode == code.ErrCodeJwtTokenIsExpired {
		api.UnauthorizedResult(ctx, errCode.Int32(), nil, errCode.String())
	} else {
		api.Result(ctx, errCode.Int32(), nil, errCode.String())
	}
}

// UnauthorizedResult method    注入401请求返包
func (api *CodeApi) UnauthorizedResult(ctx *gin.Context, code int32, data interface{}, msg string) {
	api.Code = code
	api.Message = msg
	api.Data = data
	api.NowTime = time.Now().Unix()
	if useTime := api.getUserTime(ctx); len(useTime) != 0 {
		api.UseTime = api.getUserTime(ctx)
	}
	ctx.JSON(http.StatusUnauthorized, api)
}

// Result method    注入请求返包
func (api *CodeApi) Result(ctx *gin.Context, code int32, data interface{}, msg string) {
	api.Code = code
	api.Message = msg
	api.Data = data
	api.NowTime = time.Now().Unix()
	if useTime := api.getUserTime(ctx); len(useTime) != 0 {
		api.UseTime = api.getUserTime(ctx)
	}
	ctx.JSON(http.StatusOK, api)
}

// getUserTime method    获取请求统计时间
func (api *CodeApi) getUserTime(ctx *gin.Context) string {
	startTime, _ := ctx.Get("requestStartTime")
	stopTime := time.Now().UnixMicro()
	if startTime == nil {
		return ""
	} else {
		return fmt.Sprintf("%.6f", float64(stopTime-utils.InterfaceToInt64(startTime))/1000000)
	}
}
