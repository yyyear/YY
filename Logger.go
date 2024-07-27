package YY

import (
	"encoding/json"
	"fmt"
	"github.com/robfig/cron"
	"github.com/rs/zerolog"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var Logger zerolog.Logger
var cronNew *cron.Cron

func init() {

	logger()
	i := 0
	cronNew = cron.New()
	spec := "0 0 * * * ?"
	err1 := cronNew.AddFunc(spec, func() {
		i++
		fmt.Println("cron times : ", i)
		logger()
	})
	if err1 != nil {
		fmt.Errorf("AddFunc error : %v", err1)
		return
	}
	cronNew.Start()
}

func logger() {
	timeFormat := "2006-01-02 15:04:05"
	zerolog.TimeFieldFormat = timeFormat
	now := time.Now()
	logDir := "./run_log/" + now.Format("2006-01-02")
	if !DEBUG {
		// 创建log目录
		err := os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			fmt.Println("Mkdir failed, err:", err)
			return
		}
	}

	fileName := logDir + "/" + strconv.Itoa(now.Hour()) + ".log"

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: timeFormat}
	consoleWriter.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	consoleWriter.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	consoleWriter.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	consoleWriter.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s;", i)
	}

	multi := zerolog.MultiLevelWriter(consoleWriter)
	if !DEBUG {
		logFile, _ := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		log.Println(logFile, fileName)
		multi = zerolog.MultiLevelWriter(consoleWriter, logFile)
	}

	Logger = zerolog.New(multi).With().Timestamp().Caller().Logger()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

func Relese() {
	cronNew.Stop()
}

func Info(message ...interface{}) {
	Logger.Info().CallerSkipFrame(1).Msg(ExpandText(message))
}
func Debug(message ...interface{}) {
	Logger.Debug().CallerSkipFrame(1).Msg(ExpandText(message))
}
func Error(message ...interface{}) {
	Logger.Error().CallerSkipFrame(1).Msg(ExpandText(message))
}

func ExpandArrayText(msg []interface{}) string {
	result := ""
	for _, v := range msg {
		result = result + ExpandText(v)
	}
	return result
}

func ExpandText(message interface{}) string {
	result := ""
	value := message

	var key string
	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case time.Time:
		t, _ := value.(time.Time)
		key = t.String()
		// 2022-11-23 11:29:07 +0800 CST  这类格式把尾巴去掉
		key = strings.Replace(key, " +0800 CST", "", 1)
		key = strings.Replace(key, " +0000 UTC", "", 1)
	case []byte:
		key = string(value.([]byte))
	case error:
		key = value.(error).Error()
	case []interface{}:
		key = ExpandArrayText(value.([]interface{}))
	default:
		if value == nil {
			key = " nil "
		} else {
			newValue, _ := json.Marshal(value)
			key = string(newValue)
		}

	}
	result = result + key
	return result
}
