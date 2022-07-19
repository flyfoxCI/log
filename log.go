package log

import (
	"github.com/orandin/lumberjackrus"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
)

type LogCfg struct {
	Levels     string `yaml:"levels"`     // 日志等级
	LogDir     string `yaml:"logDir"`     // 日志分类id
	MaxSize    int    `yaml:"maxSize"`    //文件大小MB
	MaxAge     int    `yaml:"maxAge"`     //文件保存时间天
	MaxBackups int    `yaml:"maxBackups"` //文件保存数量
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
func initLog(cfg *LogCfg) *log.Entry {
	logDir := cfg.LogDir
	if exists, _ := PathExists(logDir); !exists {
		err := os.Mkdir(logDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	txtFormatter := log.TextFormatter{
		DisableColors:   true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	}
	level, _ := log.ParseLevel(cfg.Levels)
	log.SetLevel(level)
	log.SetFormatter(&txtFormatter)
	maxBackups := cfg.MaxBackups
	if maxBackups == 0 {
		maxBackups = 3
	}
	maxSize := cfg.MaxSize
	if maxSize == 0 {
		maxSize = 100
	}
	maxAge := cfg.MaxAge
	if maxAge == 0 {
		maxAge = 3
	}
	hook, err := lumberjackrus.NewHook(
		&lumberjackrus.LogFile{
			Filename:   filepath.Join(logDir, "general.log"),
			MaxSize:    maxSize,
			MaxAge:     maxAge,
			MaxBackups: maxBackups,
			LocalTime:  true,
			Compress:   false,
		},
		log.DebugLevel,
		&txtFormatter,
		&lumberjackrus.LogFileOpts{
			log.DebugLevel: &lumberjackrus.LogFile{
				Filename:   filepath.Join(logDir, "debug.log"),
				MaxSize:    maxSize,
				MaxAge:     maxAge,
				MaxBackups: maxBackups,
				LocalTime:  true,
				Compress:   false,
			},
			log.InfoLevel: &lumberjackrus.LogFile{
				Filename:   filepath.Join(logDir, "info.log"),
				MaxSize:    maxSize,
				MaxAge:     maxAge,
				MaxBackups: maxBackups,
				LocalTime:  true,
				Compress:   false,
			},
			log.WarnLevel: &lumberjackrus.LogFile{
				Filename:   filepath.Join(logDir, "warn.log"),
				MaxSize:    maxSize,
				MaxAge:     maxAge,
				MaxBackups: maxBackups,
				LocalTime:  true,
				Compress:   false,
			},
			log.ErrorLevel: &lumberjackrus.LogFile{
				Filename:   filepath.Join(logDir, "error.log"),
				MaxSize:    maxSize,
				MaxAge:     maxAge,
				MaxBackups: maxBackups,
				LocalTime:  true,
				Compress:   false,
			},
			log.FatalLevel: &lumberjackrus.LogFile{
				Filename:   filepath.Join(logDir, "fatal.log"),
				MaxSize:    maxSize,
				MaxAge:     maxAge,
				MaxBackups: maxBackups,
				LocalTime:  true,
				Compress:   false,
			},
		})
	if err != nil {
		panic(err)
	}
	log.AddHook(hook)
	log.SetOutput(ioutil.Discard)
	return log.WithFields(log.Fields{"mod": "main"})
}

func getLogger() *log.Entry {
	f, err := ioutil.ReadFile("config/log.yaml")
	if err != nil {
		return nil
	}

	var logCfg *LogCfg
	err = yaml.Unmarshal(f, &logCfg)
	if err != nil {
		return nil
	}
	logger := initLog(logCfg)
	return logger
}

var Logger = getLogger()
