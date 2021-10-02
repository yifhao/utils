package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync/atomic"

	colorable "github.com/mattn/go-colorable"

	"github.com/logrusorgru/aurora"
)

// MultiLogWriter 多端写日志类
type MultiLogWriter struct {
	writers []io.Writer
	io.Writer
}

var emptyMLogWriter = MultiLogWriter{
	Writer: io.MultiWriter(),
}
var multiLoggerWriter atomic.Value
var multiLogger *log.Logger
var colorLogger = log.New(colorable.NewColorableStdout(), "", log.LstdFlags)

func init() {
	multiLoggerWriter.Store(&emptyMLogWriter)
	multiLogger = log.New(&emptyMLogWriter, "", log.LstdFlags)
	log.SetOutput(io.MultiWriter(os.Stdout, &emptyMLogWriter))
}

func getMLogWriter() *MultiLogWriter {
	inner := multiLoggerWriter.Load()
	if inner == nil {
		return &emptyMLogWriter
	}
	writer, ok := inner.(*MultiLogWriter)
	if !ok {
		return &emptyMLogWriter
	}
	return writer
}

// AddWriter 添加日志输出端
func AddWriter(wn io.Writer) {
	originalMLogWriter := getMLogWriter()

	// copy on write
	var newMLogWriter MultiLogWriter
	for i := range originalMLogWriter.writers {
		newMLogWriter.writers = append(newMLogWriter.writers, originalMLogWriter.writers[i])
	}
	newMLogWriter.writers = append(newMLogWriter.writers, wn)

	newMLogWriter.Writer = io.MultiWriter(originalMLogWriter.writers...)

	multiLoggerWriter.Store(&newMLogWriter)

	multiLogger.SetOutput(newMLogWriter)
	log.SetOutput(io.MultiWriter(os.Stdout, newMLogWriter))
}

// MayBeError 优雅错误判断加日志辅助函数
func MayBeError(info error) (hasError bool) {
	if hasError = info != nil; hasError {
		Print(aurora.Red(info))
	}
	return
}
func getNoColor(v ...interface{}) (noColor []interface{}) {
	noColor = append(noColor, v...)
	for i, value := range v {
		if vv, ok := value.(aurora.Value); ok {
			noColor[i] = vv.Value()
		}
	}
	return
}

// Print 带颜色识别
func Print(v ...interface{}) {
	noColor := getNoColor(v...)
	colorLogger.Output(2, fmt.Sprint(v...))
	multiLogger.Output(2, fmt.Sprint(noColor...))
}

// Printf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	noColor := getNoColor(v...)
	colorLogger.Output(2, fmt.Sprintf(format, v...))
	multiLogger.Output(2, fmt.Sprintf(format, noColor...))
}

// Println calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Println.
func Println(v ...interface{}) {
	noColor := getNoColor(v...)
	colorLogger.Output(2, fmt.Sprintln(v...))
	multiLogger.Output(2, fmt.Sprintln(noColor...))
}
