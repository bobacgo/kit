package r

// https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Status
// 自定义业务状态码，参考http协议的状态码 +1000

// 信息响应 (100–199)
// 成功响应 (200–299)
// 重定向消息 (300–399)
// 客户端错误响应 (400–499)
// 服务端错误响应 (500–599)

// StatusCode 响应状态
type StatusCode uint16

const (
	Ok               StatusCode = 0
	ParameterIllegal StatusCode = 4000
	ParameterInvalid StatusCode = 4001
	TokenInvalid     StatusCode = 4010
	TokenMission     StatusCode = 4011
	Forbidden        StatusCode = 4030
	Gone             StatusCode = 4100
	Internal         StatusCode = 5000
	// ....
)

// Status {code, msg}
var Status = map[StatusCode]string{
	0:    "操作成功",
	4000: "参数解析失败",
	4001: "参数校验不通过",
	4010: "无效的Token",
	4011: "Token缺失",
	4030: "权限不足",
	4050: "请求方法不支持",
	4100: "值不存在",
	5000: "系统内部错误",
	5001: "操作失败",
}