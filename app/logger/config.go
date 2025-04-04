package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/pflag"
)

const (
	timeFormatDefault     = "2006-01-02 15:04:05.000"
	filepathDefault       = "." + string(os.PathSeparator) + "logs"
	filenameDefault       = "server"
	filenameSuffixDefault = "2006-01-02-150405"
	fileExtensionDefault  = "log"
)

type Config struct {
	Level LogLevel `mapstructure:"level" validate:"omitempty,oneof=debug info error"`
	// 监听日志级别变化
	LevelCh    chan LogLevel `mapstructure:"-" json:"-" yaml:"-"`
	TimeFormat string        `mapstructure:"timeFormat" yaml:"timeFormat" default:"2006-01-02 15:04:05"`

	// 完整的文件路径名
	Filepath        string `mapstructure:"filepath"`
	Filename        string `mapstructure:"filename" default:"server"`
	FilenameSuffix  string `mapstructure:"filenameSuffix" yaml:"filenameSuffix" default:"2006-01-02-150405"`
	FileExtension   string `mapstructure:"fileExtension" yaml:"fileExtension" default:"log"`
	FileJsonEncoder bool   `mapstructure:"fileJsonEncoder" yaml:"fileJsonEncoder"`

	FileMaxSize  uint16 `mapstructure:"fileSizeMax" yaml:"fileSizeMax" default:"10"`     // 单位是MB 默认值是 10MB
	FileMaxAge   uint16 `mapstructure:"fileAgeMax" yaml:"fileAgeMax" default:"180"`      // 留存天数
	FileCompress bool   `mapstructure:"fileCompress" yaml:"fileCompress" default:"true"` // 是否归档压缩
}

func (c *Config) SetLevel(level LogLevel) {
	c.LevelCh <- level
}

// xxx/logs/xxx-service-2006-01-01-150405.log
func (c *Config) makeFilename() string {
	if c.FilenameSuffix != "" {
		c.FilenameSuffix = "-" + time.Now().Format(c.FilenameSuffix)
	}
	return filepath.Join(c.Filepath, fmt.Sprintf("%s%s.%s", c.Filename, c.FilenameSuffix, c.FileExtension))
}

type Option func(*Config)

func NewConfig(opts ...Option) Config {
	conf := Config{
		TimeFormat:     timeFormatDefault,
		Filepath:       filepathDefault,
		Filename:       filenameDefault,
		FilenameSuffix: filenameSuffixDefault,
		FileExtension:  fileExtensionDefault,
		FileMaxSize:    10,
		FileMaxAge:     6 * 30,
		FileCompress:   true,
		LevelCh:        make(chan LogLevel, 1),
	}
	for _, opt := range opts {
		opt(&conf)
	}
	return conf
}

// WithLevel 日志级别
func WithLevel(level LogLevel) Option {
	return func(o *Config) {
		o.Level = level
		o.LevelCh <- level
	}
}

// WithTimeFormat 日志时间格式
func WithTimeFormat(timeFormat string) Option {
	return func(o *Config) {
		o.TimeFormat = timeFormat
	}
}

// WithFilepath 文件目录路径
func WithFilepath(filepath string) Option {
	return func(o *Config) {
		o.Filepath = filepath
	}
}

// WithFilename 文件名(文件前缀), 随机部分 main-service-2023-11-04
func WithFilename(filename string) Option {
	return func(o *Config) {
		if filename != "" {
			o.Filename = filename
		}
	}
}

// WithFilenameSuffix 文件后缀名, 随机部分 main-service-2023-11-04
func WithFilenameSuffix(filenameSuffix string) Option {
	return func(o *Config) {
		o.FilenameSuffix = filenameSuffix
	}
}

// WithFileExtension 文件扩展名 (e.g log、txt)
func WithFileExtension(fileExtension string) Option {
	return func(o *Config) {
		if fileExtension != "" {
			o.FileExtension = fileExtension
		}
	}
}

// WithFileJsonEncoder 输出到文件侧是否启用json格式编码
func WithFileJsonEncoder(isJsonEncoder bool) Option {
	return func(o *Config) {
		o.FileJsonEncoder = isJsonEncoder
	}
}

// WithFileMaxSize 文件最大多少MB就分割
func WithFileMaxSize(maxSize uint16) Option {
	return func(o *Config) {
		o.FileMaxSize = maxSize
	}
}

// WithFileMaxAge 文件保留时长
func WithFileMaxAge(maxAge uint16) Option {
	return func(o *Config) {
		o.FileMaxAge = maxAge
	}
}

// WithFileCompress 是否归档压缩
func WithFileCompress(compress bool) Option {
	return func(o *Config) {
		o.FileCompress = compress
	}
}

func (c *Config) Validate() []error {
	// TODO valid config data
	return nil
}

func (c *Config) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&c.Filename, "log-filename", c.Filename, "log filename")
}

func New(appName string, logCfg Config) Config {
	if logCfg.Level == "" {
		logCfg.Level = LogLevel_Info
	}
	opts := []Option{
		WithLevel(logCfg.Level),
	}
	if logCfg.TimeFormat != "" {
		opts = append(opts, WithTimeFormat(logCfg.TimeFormat))
	}
	if appName != "" {
		opts = append(opts, WithFilename(appName))
	}
	if logCfg.Filepath != "" {
		opts = append(opts, WithFilepath(logCfg.Filepath))
	}
	if logCfg.FilenameSuffix != "" {
		opts = append(opts, WithFilenameSuffix(logCfg.FilenameSuffix))
	}
	if logCfg.FileExtension != "" {
		opts = append(opts, WithFileExtension(logCfg.FileExtension))
	}
	if logCfg.FileMaxSize > 0 {
		opts = append(opts, WithFileMaxSize(logCfg.FileMaxSize))
	}
	if logCfg.FileMaxAge > 0 {
		opts = append(opts, WithFileMaxAge(logCfg.FileMaxAge))
	}
	if logCfg.FileJsonEncoder {
		opts = append(opts, WithFileJsonEncoder(logCfg.FileJsonEncoder))
	}
	if logCfg.FileCompress {
		opts = append(opts, WithFileCompress(logCfg.FileCompress))
	}

	cfg := NewConfig(opts...)
	// 初始化日志配置
	InitZapLogger(cfg)
	return cfg
}
