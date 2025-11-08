package logger

import (
	"context"
	"os"
	"path/filepath"

	"github.com/narcissus1949/narcissus-blog/internal/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger         *zap.Logger
	DefaultLogPath = filepath.Join(utils.GetRootDir(), "log")
)

type LogConfig struct {
	LogLevel          string `json:"logLevel,omitempty" yaml:"logLevel,omitempty"`                   // 日志打印级别 debug  info  warning  error
	LogFormat         string `json:"logFormat,omitempty" yaml:"logFormat,omitempty"`                 // 输出日志格式	logfmt, json
	LogPath           string `json:"logPath,omitempty" yaml:"logPath,omitempty"`                     // 输出日志文件路径
	LogFileName       string `json:"logFileName,omitempty" yaml:"logFileName,omitempty"`             // 输出日志文件名称
	LogFileMaxSize    int    `json:"logFileMaxSize,omitempty" yaml:"logFileMaxSize,omitempty"`       // 【日志分割】单个日志文件最多存储量 单位(mb)
	LogFileMaxBackups int    `json:"logFileMaxBackups,omitempty" yaml:"logFileMaxBackups,omitempty"` // 【日志分割】日志备份文件最多数量
	LogMaxAge         int    `json:"logMaxAge,omitempty" yaml:"logMaxAge,omitempty"`                 // 日志保留时间，单位: 天 (day)
	LogCompress       bool   `json:"logCompress,omitempty" yaml:"logCompress,omitempty"`             // 是否压缩日志
	LogStdout         bool   `json:"logStdout,omitempty" yaml:"logStdout,omitempty"`                 // 是否输出到控制台
}

func NewDefaultCfg() LogConfig {
	return LogConfig{
		LogLevel:          "info",
		LogFormat:         "logfmt",
		LogPath:           DefaultLogPath,
		LogFileName:       "blog.log",
		LogFileMaxSize:    1024,
		LogFileMaxBackups: 30,
		LogMaxAge:         30,
		LogCompress:       false,
		LogStdout:         true,
	}
}

func MustInit(cfg LogConfig) {
	logLevel := map[string]zapcore.Level{
		"debug": zapcore.DebugLevel,
		"info":  zapcore.InfoLevel,
		"warn":  zapcore.WarnLevel,
		"error": zapcore.ErrorLevel,
	}
	writeSyncer, err := getLogWriter(cfg)
	if err != nil {
		panic(err.Error())
	}
	encoder := getEncoder(cfg)
	level, ok := logLevel[cfg.LogLevel]
	if !ok {
		panic("log level not support")
	}
	core := zapcore.NewCore(encoder, writeSyncer, level)
	logger = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger) // 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
}

// getEncoder 编码器(如何写入日志)
func getEncoder(conf LogConfig) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   // log 时间格式 例如: 2021-09-11t20:05:54.852+0800
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 输出level序列化为全大写字符串，如 INFO DEBUG ERROR
	if conf.LogFormat == "json" {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// getLogWriter 获取日志输出方式  日志文件 控制台
func getLogWriter(conf LogConfig) (zapcore.WriteSyncer, error) {
	if exist := IsExist(conf.LogPath); !exist {
		if conf.LogPath == "" {
			conf.LogPath = DefaultLogPath
		}
		if err := os.MkdirAll(conf.LogPath, os.ModePerm); err != nil {
			return nil, err
		}
	}

	// 日志文件 与 日志切割 配置
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filepath.Join(conf.LogPath, conf.LogFileName), // 日志文件路径
		MaxSize:    conf.LogFileMaxSize,                           // 单个日志文件最大多少 mb
		MaxBackups: conf.LogFileMaxBackups,                        // 日志备份数量
		MaxAge:     conf.LogMaxAge,                                // 日志最长保留时间
		Compress:   conf.LogCompress,                              // 是否压缩日志
	}
	if conf.LogStdout {
		// 日志同时输出到控制台和日志文件中
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(lumberJackLogger), zapcore.AddSync(os.Stdout)), nil
	} else {
		// 日志只输出到日志文件
		return zapcore.AddSync(lumberJackLogger), nil
	}
}

// IsExist 判断文件或者目录是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// ---- Context helpers ----
// 使用私有 key 避免与其他 context 值冲突
type ctxKey struct{}

var loggerCtxKey = ctxKey{}

// FromContext 从 context 中获取 *zap.Logger
// 如果未设置，则回退到全局 logger（zap.L()）
func FromContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return zap.L()
	}
	if v := ctx.Value(loggerCtxKey); v != nil {
		if l, ok := v.(*zap.Logger); ok && l != nil {
			return l
		}
	}
	return zap.L()
}

// ToContext 将 *zap.Logger 放入 context，返回新的 context
func ToContext(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, l)
}
