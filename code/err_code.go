package code

type ErrCode int32

type IErrCode interface {
	Int32() int32
	Int64() int64
	UInt32() uint32
	String() string
}

const (
	Success   ErrCode = 200 // 成功
	NoErrCode ErrCode = 0   // 无错误码
)

// 协议错误码
const (
	ErrCodeOfflineReasonHandlerPanic ErrCode = 100 // 离线处理异常
	ErrCodeServerClosePushMsgNo      ErrCode = 101 // 离线处理没有消息推送
	ErrCodeAlreadyLogin              ErrCode = 102 // 已经登录
	EErrCodeRepeatedLogin            ErrCode = 103 // 已经登录
	ErrCodeFileNotExist              ErrCode = 404 // 资源不存在

	ErrCodeInvalidParams ErrCode = 500 // 请求参数错误
	ErrCodeJwtTokenErr   ErrCode = 501 // token错误
	ErrCodeConnRefuse    ErrCode = 503 // 连接被拒绝
	ErrCodeNetAbnormal   ErrCode = 504 // 网络异常
	ErrCodeTimeOut       ErrCode = 505 // 请求超时
	ErrCodeRedisWriteErr ErrCode = 506 // Redis写入错误

	ErrCodeJwtTokenIsExpired       ErrCode = 600 // jwt的token过期
	ErrCodeJwtTokenNotActiveYet    ErrCode = 601 // jwt的没有启用
	ErrCodeJwtNotEvenAToken        ErrCode = 602 // 没有携带jwt的token
	ErrCodeJwtTokenNotInvalid      ErrCode = 603 // jwt的token不正确
	ErrCodeAuthorizationRefreshErr ErrCode = 604 // jwt的token刷新错误
	ErrCodeJwtGenerateErr          ErrCode = 605 // jwt生成失败

	ErrCodeCasbinNotActiveYet    ErrCode = 800 // casbin没有启用
	ErrCodeCasbinNotPermissions  ErrCode = 801 // casbin没有权限
	ErrCodeDeleteCasbinGlobalErr ErrCode = 802 // casbin删除全局权限失败
)

// 业务错误码
const (
	ErrCodeGenerateUUidErr               ErrCode = 1000 // 生成uuid错误
	ErrCodeGenerateAccountErr            ErrCode = 1001 // 生成账户失败
	ErrCodeUuidGetDBErr                  ErrCode = 1002 // 从数据库获取uuid错误
	ErrCodeAccountNotExist               ErrCode = 1003 // 账户不存在
	ErrCodeUserNameExist                 ErrCode = 1004 // 账户名已存在
	ErrCodePasswordErr                   ErrCode = 1005 // 密码错误
	ErrCodeNotBanSelf                    ErrCode = 1006 // 不能禁用自己
	ErrCodeRoleBan                       ErrCode = 1007 // 账户被禁用
	ErrCodeSmsCodeInvalid                ErrCode = 1008 // 短信验证码错误
	ErrCodeRoleNotExist                  ErrCode = 1009 // 角色不存在
	ErrCodeRolePermissionErr             ErrCode = 1010 // 角色权限不够
	ErrCodeCasbinIdenticalAdditionFailed ErrCode = 1011 // casbin没找到相同api

	ErrCodeNotFoundAESKey   ErrCode = 1200 // 没找到aes密钥
	ErrCodeAESDecodingError ErrCode = 1201 // aes解码错误

	ErrCodeNotFoundMenusGlobal  ErrCode = 1300 // 没找路由表
	ErrCodeNotGetAllMenusGlobal ErrCode = 1301 // 没有导出路由
	ErrCodeDeleteMenusGlobal    ErrCode = 1302 // 删除路由失败

	ErrCodeNotFoundUserShare ErrCode = 1400 // 没找到用户表
	ErrCodeCreateUserShare   ErrCode = 1401 // 生成用户表失败
	ErrCodeDeleteUserShare   ErrCode = 1402 //  删除用户表失败
	ErrCodePersistUserShare  ErrCode = 1403 // 持久化用户表失败

	ErrCodeNotFoundAccountGlobal ErrCode = 1500 // 没找到账户表

	ErrCodeNotFountRoleGlobal  ErrCode = 1600 // 没找到角色表
	ErrCodeRoleCodeExist       ErrCode = 1601 // 角色编码已存在
	ErrCodeNotFoundRoleGlobals ErrCode = 1602 // 没找到角色表
	ErrCodeDeleteRoleGlobalErr ErrCode = 1603 //  删除角色表失败

	ErrCodeNotFoundApisGlobals ErrCode = 1700 // 没找到接口表
	ErrCodeApiPathExist        ErrCode = 1701 // 接口路径已存在
	ErrCodeDeleteApisErr       ErrCode = 1702 // 删除接口表错误

	ErrCodeNotFoundRoleAuthGlobal ErrCode = 1800 // 没找到角色权限表
	ErrCodeNotFoundRoleApiGlobal  ErrCode = 1801 // 没找到角色接口表
	ErrCodeDeleteRoleApiGlobal    ErrCode = 1802 // 删除角色接口表错误
	ErrCodeDeleteRoleAuthGlobal   ErrCode = 1803 // 删除角色权限表错误

	ErrCodeNotFoundRecordsGlobals ErrCode = 1900 // 没找到记录表
	ErrCodeDeleteRecordsGlobals   ErrCode = 1901 // 删除记录表错误

	ErrCodeNotFoundDictTypeGlobal  ErrCode = 2000 // 没找到字典类型表
	ErrCodeDictTypeGlobalExitErr   ErrCode = 2001 // 字典类型编码已存在
	ErrCodeDictTypeGlobalDeleteErr ErrCode = 2002 // 字典类型删除失败
	ErrCodeDictDataGlobalDeleteErr ErrCode = 2002 // 字典数据删除失败
	ErrCodeDictInvalidParams       ErrCode = 2003 // 请求获取字典参数错误
	ErrCodeNotFoundDictDataGlobal  ErrCode = 2004 // 没找到字典数据表

	ErrCodeAssetGlobalExist    ErrCode = 2100 // 资产名称已存在
	ErrCodeAssetGlobalNotExist ErrCode = 2101 // 资产不存在

	ErrCodeNotFoundGpuMonitor ErrCode = 2200 // 没找到gpu监控表
	ErrCodeGpuMonitorMarshal  ErrCode = 2201 // gpu监控表序列化错误

	ErrCodeDBSyncErr ErrCode = 9000 // 数据库同步错误
)

func (e ErrCode) Int32() int32 {
	return int32(e)
}

func (e ErrCode) UInt32() uint32 {
	return uint32(e)
}

func (e ErrCode) Int64() int64 {
	return int64(e)
}
