package codes

// https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Status#服务端错误响应
/*
	信息响应 (100–199)
	成功响应 (200–299)
	重定向消息 (300–399)
	客户端错误响应 (400–499)
	服务端错误响应 (500–599)
*/
// Code 响应状态
type Code = uint32

const (
	OK                  Code = 200
	BadRequest          Code = 400
	TokenInvalid        Code = 4010
	TokenMission        Code = 4011
	InternalServerError Code = 500
	// ....
	RecordRepeat Code = 5000
	DateBusy     Code = 5001
)
