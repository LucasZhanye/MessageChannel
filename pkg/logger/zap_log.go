package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type zapLog struct {
	log *zap.Logger
}

// 默认输出日志文件路径
const DefaultLogPath = "./runtime/log"
const DefaultLogFileName = "debug.log"

func InitZapLog(conf Config) (Log, error) {
	zapLog := &zapLog{}

	logLevel := map[string]zapcore.Level{
		"debug": zap.DebugLevel,
		"info":  zap.InfoLevel,
		"warn":  zap.WarnLevel,
		"error": zap.ErrorLevel,
		"fatal": zap.FatalLevel,
		"panic": zap.PanicLevel,
	}

	// 日志文件配置，文件位置和切割
	writeSyncer, err := getLogWriter(conf)
	if err != nil {
		return nil, err
	}
	// 获取日志输出编码
	encoder := getEncoder(conf)
	// 日志打印级别
	level, ok := logLevel[conf.Level]
	if !ok {
		level = logLevel["info"]
	}

	core := zapcore.NewCore(encoder, writeSyncer, level)
	// zap.AddCaller() 记录文件名和行号
	zapLog.log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	// 1. zap.ReplaceGlobals 函数将当前初始化的 logger 替换到全局的 logger,
	// 2. 使用 logger 的时候 直接通过 zap.S().Debugf("xxx") or zap.L().Debug("xxx")
	// 3. 使用 zap.S() 和 zap.L() 提供全局锁，保证一个全局的安全访问logger的方式
	zap.ReplaceGlobals(zapLog.log)
	return zapLog, nil
}

// getEncoder 编码器
func getEncoder(conf Config) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	// log时间格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 输出level序列化为全大写字符
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// json格式输出时时间的key
	encoderConfig.TimeKey = "time"
	// 时间格式
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")

	// 日志输出格式
	if conf.Format == "json" {
		return zapcore.NewJSONEncoder(encoderConfig)
	}

	return zapcore.NewConsoleEncoder(encoderConfig)
}

// getLogWriter
func getLogWriter(conf Config) (zapcore.WriteSyncer, error) {
	// 判断日志路径是否存在，不存在则创建
	if exist := isExist(conf.Path); !exist {
		if conf.Path == "" {
			conf.Path = DefaultLogPath
		}
		if err := os.MkdirAll(conf.Path, os.ModePerm); err != nil {
			conf.Path = DefaultLogPath
			if err := os.MkdirAll(conf.Path, os.ModePerm); err != nil {
				return nil, err
			}
		}
	}

	if conf.FileName == "" {
		conf.FileName = DefaultLogFileName
	}

	// 日志文件与日志分割配置
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filepath.Join(conf.Path, conf.FileName), // 日志文件路径
		MaxSize:    conf.FileMaxSize,
		MaxBackups: conf.FileMaxBackup,
		MaxAge:     conf.MaxAge,
		Compress:   conf.Compress,
		LocalTime:  true, // 使用本地时间
	}

	if conf.StdOut {
		// 日志同时输出到控制台和日志文件中
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(lumberJackLogger), zapcore.AddSync(os.Stdout)), nil
	} else {
		return zapcore.AddSync(lumberJackLogger), nil
	}
}

// isExist 判断文件或者目录是否存在
func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func (l *zapLog) Info(format string, args ...interface{}) {
	zap.S().Infof(format, args...)
}

func (l *zapLog) Debug(format string, args ...interface{}) {
	zap.S().Debugf(format, args...)
}

func (l *zapLog) Error(format string, args ...interface{}) {
	zap.S().Errorf(format, args...)
}

func (l *zapLog) Warn(format string, args ...interface{}) {
	zap.S().Warnf(format, args...)
}

func (l *zapLog) Fatal(format string, args ...interface{}) {
	zap.S().Fatalf(format, args...)
}

func (l *zapLog) Panic(format string, args ...interface{}) {
	zap.S().Panicf(format, args...)
}

func (l zapLog) Println(args ...interface{}) {
	zap.S().Infoln(args)
}

func (l zapLog) Printf(format string, args ...interface{}) {
	zap.S().Infof(format, args)
}
