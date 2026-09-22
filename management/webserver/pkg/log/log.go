package log

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"

	"chaitin.cn/dev/go/log"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/config"
	"chaitin.cn/patronus/safeline-2/management/webserver/utils"
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

		// O_NOFOLLOW closes the window between Stat and OpenFile: a symlink swapped
		// in for the log path is refused instead of followed.
		fileFlag := os.O_WRONLY | os.O_APPEND | os.O_SYNC | syscall.O_NOFOLLOW
		if !exist {
			if err := utils.EnsureFileDir(config.GlobalConfig.Log.Output); err != nil {
				return err
			}
			fileFlag = fileFlag | os.O_CREATE
		}

		// os.ModePerm (0777) let every local user rewrite the audit log. The
		// file is created with an owner/group readable mode instead.
		if fp, err := os.OpenFile(config.GlobalConfig.Log.Output, fileFlag, logFileMode); err != nil {
			return fmt.Errorf("failed to open log file: %s", err.Error())
		} else {
			log.SetOutput(log.AllLoggers, log.NewLockOutput(fp))
		}

		// The mode above only applies when the file is created, so a log file
		// that an older release left world writable is tightened here.
		tightenLogFileMode(config.GlobalConfig.Log.Output)
	}

	// hook
	log.AddHook(log.AllLoggers, NewRuntimeHook())
	log.AddHook(log.AllLoggers, log.NewErrorStackHook(true))
	// level
	LoadLogLevel()

	return nil
}

// logFileMode is the permission the log file is created with. Logs are audit
// material: they must not be writable by other local users.
const logFileMode os.FileMode = 0640

// looseLogFileMask matches the group and world write bits.
const looseLogFileMask os.FileMode = 0o022

// tightenLogFileMode removes the group and world write bits from an existing
// log file, so that a log file created by an older release cannot stay
// tamperable by other local users.
func tightenLogFileMode(path string) {
	info, err := os.Stat(path)
	if err != nil || info.Mode().Perm()&looseLogFileMask == 0 {
		return
	}

	log.Warnf("Log file %s is writable by other users, setting mode to %#o", path, logFileMode)
	if err := os.Chmod(path, logFileMode); err != nil {
		log.Warnf("Cannot tighten the mode of log file %s: %s", path, err)
	}
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
