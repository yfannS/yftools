package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Level 日志级别
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

var levelNames = map[Level]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	FATAL: "FATAL",
}

// Logger 结构化日志器
type Logger struct {
	mu      sync.Mutex
	level   Level
	prefix  string
	logger  *log.Logger
	outputs []io.Writer
}

var (
	defaultLogger *Logger
	once          sync.Once
)

func init() {
	defaultLogger = NewLogger("")
}

// NewLogger 创建新 Logger，默认输出到 stdout
func NewLogger(prefix string) *Logger {
	l := &Logger{
		level:   INFO,
		prefix:  prefix,
		outputs: []io.Writer{os.Stdout},
	}
	l.logger = log.New(io.MultiWriter(l.outputs...), "", 0)
	return l
}

// SetLevel 设置日志级别
func SetLevel(level Level) {
	defaultLogger.SetLevel(level)
}

func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// SetPrefix 设置日志前缀（模块名）
func SetPrefix(prefix string) {
	defaultLogger.SetPrefix(prefix)
}

func (l *Logger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prefix = prefix
}

// AddFileOutput 添加文件输出
func AddFileOutput(path string) error {
	return defaultLogger.AddFileOutput(path)
}

func (l *Logger) AddFileOutput(path string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create log dir: %w", err)
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("open log file: %w", err)
	}

	l.outputs = append(l.outputs, f)
	l.logger.SetOutput(io.MultiWriter(l.outputs...))
	return nil
}

// SetOutput 设置自定义输出（替换默认 stdout）
func SetOutput(w io.Writer) {
	defaultLogger.SetOutput(w)
}

func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.outputs = []io.Writer{w}
	l.logger.SetOutput(io.MultiWriter(l.outputs...))
}

// ParseLevel 将字符串解析为 Level
func ParseLevel(s string) Level {
	switch s {
	case "DEBUG", "debug":
		return DEBUG
	case "INFO", "info":
		return INFO
	case "WARN", "warn":
		return WARN
	case "ERROR", "error":
		return ERROR
	case "FATAL", "fatal":
		return FATAL
	default:
		return INFO
	}
}

func (l *Logger) log(level Level, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now().Format("2006-01-02 15:04:05.000")
	msg := fmt.Sprintf(format, args...)

	prefix := l.prefix
	if prefix != "" {
		prefix = "[" + prefix + "]"
	}

	output := fmt.Sprintf("%s %s%s %s", now, levelNames[level], prefix, msg)

	switch level {
	case FATAL:
		l.logger.Fatal(output)
	default:
		l.logger.Println(output)
	}
}

// ====== 全局函数（使用 defaultLogger） ======

func Debug(format string, args ...interface{}) {
	defaultLogger.log(DEBUG, format, args...)
}

func Info(format string, args ...interface{}) {
	defaultLogger.log(INFO, format, args...)
}

func Warn(format string, args ...interface{}) {
	defaultLogger.log(WARN, format, args...)
}

func Error(format string, args ...interface{}) {
	defaultLogger.log(ERROR, format, args...)
}

func Fatal(format string, args ...interface{}) {
	defaultLogger.log(FATAL, format, args...)
}

// ====== Logger 实例方法 ======

func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
}
