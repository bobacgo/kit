package logger

// --------------- log ----------------

// Config config model
//
// 日志文件名 xxx/logs/${App.Service}-2006-01-01-150405.log
type Config struct {
	Level       string // 日志级别 默认值是 info
	FileSizeMax uint16 `yaml:"fileSizeMax"`             // 单位是MB 默认值是 10MB
	FileAgeMax  uint16 `yaml:"fileAgeMax"`              // 留存天数
	DirPath     string `validator:"dir" yaml:"dirPath"` // 日志文件夹路径 默认 ./logs
}