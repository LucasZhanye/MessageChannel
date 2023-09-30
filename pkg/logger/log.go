package logger

type Config struct {
	Level         string `json:"level" yaml:"level"`                 // 日志打印级别：debug，info，warning，error
	Format        string `json:"format" yaml:"format"`               // 日志输出格式： json, logfmt
	Path          string `json:"path" yaml:"path"`                   // 日志文件路径
	FileName      string `json:"fileName" yaml:"fileName"`           // 日志文件名
	FileMaxSize   int    `json:"fileMaxSize" yaml:"fileMaxSize"`     // 日志分割，单个日志文件最多存储量，单位 mb
	FileMaxBackup int    `json:"fileMaxBackup" yaml:"fileMaxBackup"` // 日志分割，日志备份文件最多数量
	MaxAge        int    `json:"maxAge" yaml:"maxAge"`               // 日志保留时间，单位：天
	Compress      bool   `json:"compress" yaml:"compress"`           // 是否压缩日志
	StdOut        bool   `json:"stdout" yaml:"stdout"`               // 是否输出到控制台
}

type Log interface {
	Println(args ...interface{})
	Printf(format string, args ...interface{})
	Info(format string, args ...interface{})
	Debug(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
	Panic(format string, args ...interface{})
}
