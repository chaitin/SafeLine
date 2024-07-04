package log

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"

	"chaitin.cn/dev/go/log"
	"chaitin.cn/patronus/safeline-2/management/tcontrollerd/pkg/config"
	"chaitin.cn/patronus/safeline-2/management/tcontrollerd/utils"
)

func GetLogger(name string) *log.Logger {
	return log.GetLogger(name)
}

func LoadLogLevel() {
	lv, _ := log.ParseLevel(config.GlobalConfig.Log.Level)
	log.SetLevel(log.AllLoggers, lv)
}

func SetLogFormatter() {
	// format
	formatter := new(log.TextFormatter)
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "2006/01/02 15:04:05"
	log.SetFormatter(log.AllLoggers, formatter)
}

func InitLogger() error {
	// output
	switch config.GlobalConfig.Log.Output {
	case "stdout":
		log.SetOutput(log.AllLoggers, os.Stdout)
	case "stderr":
		log.SetOutput(log.AllLoggers, os.Stderr)
	default:
		exist, err := utils.FileExist(config.GlobalConfig.Log.Output)
		if err != nil {
			return err
		}

		fileFlag := os.O_WRONLY | os.O_APPEND | os.O_SYNC
		if !exist {
			if err := utils.EnsureFileDir(config.GlobalConfig.Log.Output); err != nil {
				return err
			}
			fileFlag = fileFlag | os.O_CREATE
		}

		if fp, err := os.OpenFile(config.GlobalConfig.Log.Output, fileFlag, os.ModePerm); err != nil {
			return fmt.Errorf("failed to open log file: %s", err.Error())
		} else {
			log.SetOutput(log.AllLoggers, log.NewLockOutput(fp))
		}
	}

	// hook
	log.AddHook(log.AllLoggers, NewRuntimeHook())
	log.AddHook(log.AllLoggers, log.NewErrorStackHook(true))
	// level
	LoadLogLevel()

	return nil
}

type RuntimeHook struct{}

func (h *RuntimeHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *RuntimeHook) Fire(entry *logrus.Entry) error {
	file := "???"
	funcName := "???"
	line := 0

	pc := make([]uintptr, 64)
	// Skip runtime.Callers, self, and another call from logrus
	n := runtime.Callers(3, pc)
	if n != 0 {
		pc = pc[:n] // pass only valid pcs to runtime.CallersFrames
		frames := runtime.CallersFrames(pc)

		// Loop to get frames.
		// A fixed number of pcs can expand to an indefinite number of Frames.
		for {
			frame, more := frames.Next()
			if !strings.Contains(frame.File, "github.com/sirupsen/logrus") && !strings.Contains(frame.Function, "chaitin.cn/dev/go") {
				file = frame.File
				funcName = frame.Function
				line = frame.Line
				break
			}
			if !more {
				break
			}
		}
	}

	slices := strings.Split(file, "/")
	file = slices[len(slices)-1]

	funcName = strings.ReplaceAll(funcName, "chaitin.cn", "")

	entry.Data["file"] = file
	entry.Data["func"] = funcName
	entry.Data["line"] = line
	return nil
}

func NewRuntimeHook() *RuntimeHook {
	return &RuntimeHook{}
}
