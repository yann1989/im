package file_log

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

const (
	LOG_LEVEL_DEBUG = iota
	LOG_LEVEL_INFO
	LOG_LEVEL_WARNING
	LOG_LEVEL_ERROR
	LOG_LEVEL_CRITIC
)

const (
	defaultDepth = 10
	DEL_LOG_DAYS = 7
	DAY_SECONDS  = 24 * 3600
)

var logDir = *logD

type stLogger struct {
	m_FileDir       string
	m_FileName      string
	m_FileHandle    *os.File
	m_Level         int
	m_Depth         int
	m_nexDay        time.Time
	m_MaxLogFileNum int
	m_MaxLogDataNum int
	m_DelDay        uint
	m_mu            sync.Mutex
}

var (
	ArgsInvaild      = errors.New("args can be vaild")
	ObtainFileFail   = errors.New("obtain file failed")
	OpenFileFail     = errors.New("open file failed")
	GetLineNumFail   = errors.New("get line num faild")
	WriteLogInfoFail = errors.New("write log msg failed")
	LogFileError     = errors.New("log file path invaild")
)

func defaultNew() *stLogger {
	return &stLogger{
		m_FileDir:       *logD,
		m_FileName:      "",
		m_FileHandle:    nil,
		m_Level:         0,
		m_Depth:         defaultDepth,
		m_MaxLogFileNum: *maxFileNum,
		m_MaxLogDataNum: *maxFileCap,
		m_DelDay:        *delDay,
	}
}

func NewRealStLogger(level int) *stLogger {
	logger := defaultNew()
	logger.m_Depth = defaultDepth

	if level < LOG_LEVEL_DEBUG || level > LOG_LEVEL_CRITIC {
		fmt.Println("level is invailed")
	}
	logger.m_Level = level
	err := logger.obtainLofFile()
	if err != nil {
		fmt.Println(ObtainFileFail)
	}
	return logger
}

func (this *stLogger) SetLoggerDepth(depth int) {
	if depth > 0 {
		this.m_Depth = depth
	}
}

func (this *stLogger) obtainLofFile() error {
	fileDir := this.m_FileDir
	//文件夹为空
	if fileDir == "" {
		fmt.Println(ArgsInvaild)
		os.Exit(1)
	}

	//时间文件夹
	destFilePath := fmt.Sprintf("%s%d%d%d", logDir, time.Now().Year(), time.Now().Month(),
		time.Now().Day())
	flag, err := IsExist(destFilePath)
	if err != nil {
		fmt.Println(ArgsInvaild)
	}
	if !flag {
		_ = os.MkdirAll(destFilePath, os.ModePerm)
	}
	//文件夹存在,直接以创建的方式打开文件
	destFilePath = destFilePath + "/"
	logFilePath := fmt.Sprintf("%s%s_%d_%d%d%d%s", destFilePath, "log", 1, time.Now().Year(), time.Now().Month(),
		time.Now().Day(), ".log")

	_, fileSize := GetFileByteSize(logFilePath)
	if flag && fileSize > int64(this.m_MaxLogDataNum) {
		this.RenameTooBigFile()
	}
	fileHandle, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(OpenFileFail, err.Error())
	}

	this.m_FileHandle = fileHandle
	this.m_FileName = logFilePath
	//设置下次创建文件的时间
	time.Unix(time.Now().Unix(), 0).Format("2006-01-02")
	nextDay := time.Unix(time.Now().Unix()+(24*3600), 0)
	nextDay = time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 0, 0, 0,
		0, nextDay.Location())
	this.m_nexDay = nextDay

	return nil
}

func (this *stLogger) RenameTooBigFile() {
	destFilePath := fmt.Sprintf("%s%d%d%d", logDir, time.Now().Year(), time.Now().Month(),
		time.Now().Day())
	destFilePath += "/"
	for i := this.m_MaxLogFileNum; i > 0; i-- {
		newFileName := fmt.Sprintf("%s%s_%d_%d%d%d%s", destFilePath, "log", i, time.Now().Year(), time.Now().Month(),
			time.Now().Day(), ".log")
		oldFileName := fmt.Sprintf("%s%s_%d_%d%d%d%s", destFilePath, "log", i-1, time.Now().Year(), time.Now().Month(),
			time.Now().Day(), ".log")

		_ = os.Rename(oldFileName, newFileName)
	}
}

// 格式化写入日志文件
func (this *stLogger) FormatWriteLogMsg(level int, logMsg string) {
	this.m_mu.Lock()
	defer this.m_mu.Unlock()
	now := time.Now()
	//超时或者超过大小
	_, fileSize := GetFileByteSize(this.m_FileName)
	if now.Unix() > this.m_nexDay.Unix() ||
		int(fileSize) > this.m_MaxLogDataNum {
		err := this.obtainLofFile()
		if err != nil {
			fmt.Println(ObtainFileFail)
		}
	}
	for i := 0; i < DEL_LOG_DAYS; i++ {
		this.RemoveTimeOutLogFolder(this.m_DelDay + uint(i))
	}
	//flag := GetLoggerLevel(level)
	//
	//_, file, line, ok := runtime.Caller(this.m_Depth)
	//if ok == false {
	//	fmt.Println(GetLineNumFail)
	//}
	//name := path.Base(file)
	//times := time.Now().Format("2006-01-02 15:04:05")
	//_, err := Write(this.m_FileHandle, fmt.Sprintf("%s %s [%s:%d] %s\n", times, flag, name, line, logMsg))
	_, err := Write(this.m_FileHandle, logMsg)
	if err != nil {
		fmt.Println(WriteLogInfoFail, err.Error())
	}
}

func (this *stLogger) DEBUG(format string, args ...interface{}) {
	if LOG_LEVEL_DEBUG < this.m_Level {
		return
	}
	this.FormatWriteLogMsg(LOG_LEVEL_DEBUG, fmt.Sprintf(format, args...))
}

func (this *stLogger) INFO(format string, args ...interface{}) {
	if LOG_LEVEL_INFO < this.m_Level {
		return
	}
	this.FormatWriteLogMsg(LOG_LEVEL_INFO, fmt.Sprintf(format, args...))
}

func (this *stLogger) WARNING(format string, args ...interface{}) {
	if LOG_LEVEL_WARNING < this.m_Level {
		return
	}
	this.FormatWriteLogMsg(LOG_LEVEL_WARNING, fmt.Sprintf(format, args...))
}

func (this *stLogger) ERROR(format string, args ...interface{}) {
	if LOG_LEVEL_ERROR < this.m_Level {
		return
	}
	this.FormatWriteLogMsg(LOG_LEVEL_ERROR, fmt.Sprintf(format, args...))
}

func (this *stLogger) CRITIC(format string, args ...interface{}) {
	if LOG_LEVEL_CRITIC < this.m_Level {
		return
	}
	this.FormatWriteLogMsg(LOG_LEVEL_CRITIC, fmt.Sprintf(format, args...))
}

func (this *stLogger) RemoveTimeOutLogFolder(uiDayAgo uint) {
	timeNow := time.Now().Unix()
	timeAgo := timeNow - int64(uiDayAgo*DAY_SECONDS)
	t := time.Unix(timeAgo, 0)
	folderName := fmt.Sprintf("%s%d%d%d", logDir, t.Year(), t.Month(), t.Day())
	_ = os.RemoveAll(folderName)
}

func GetLoggerLevel(level int) string {
	switch level {
	case LOG_LEVEL_DEBUG:
		return "[DEBUG]:"
	case LOG_LEVEL_INFO:
		return "[INFO]:"
	case LOG_LEVEL_WARNING:
		return "[WARNING]:"
	case LOG_LEVEL_ERROR:
		return "[ERROR]:"
	case LOG_LEVEL_CRITIC:
		return "[CRITIC]:"
	default:
		return ""
	}
}
