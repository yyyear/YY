package YY

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var DEBUG bool = true

func Str(i int64, base int) string {

	return strconv.FormatInt(i, base)
}
func StrInt32(i int32, base int) string {

	return Str(int64(i), base)
}
func StrInt64(i uint64) string {
	return strconv.FormatUint(i, 10)
}
func Int(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func Int32(s string) int32 {
	return int32(Int(s))
}

func Int64(s string) int64 {
	return int64(Int(s))
}

// HitOrMiss 百分比命中
func HitOrMiss(ration float32) bool {
	return rand.Int31n(10000) < int32(ration*100)
}

// EmailValid 邮箱有效性验证
func EmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}

var numberString = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
var codeString = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "J", "K", "L", "M", "N", "P", "Q", "R", "S", "T", "W", "X", "Y", "Z", "a", "b", "c", "d", "e", "f", "g", "h", "j", "k", "m", "n", "p", "q", "r", "s", "t", "w", "x", "y", "z"}

func CreateNumber(bit int) string {
	var result = ""
	for i := 0; i < bit; i++ {
		result = result + numberString[rand.Int31n(10)]
	}
	return result
}

func CreateRandomID(bit int) string {
	var result = ""
	for i := 0; i < bit; i++ {

		result = result + codeString[rand.Int31n(int32(len(codeString)))]
	}
	return result
}

func Copy(original string) string {
	// 将字符串转换为字节切片
	copyBytes := make([]byte, len(original))
	// 复制原始字符串的内容到新的字节切片
	copy(copyBytes, original)
	// 将新的字节切片转换为字符串
	return string(copyBytes)
}

// AllNumberValid 全数字有效性验证
func AllNumberValid(e string, digit int) bool {
	res := fmt.Sprintf("^[0-9]{%d}$", digit)
	numRegex := regexp.MustCompile(res)
	return numRegex.MatchString(e)
}

// YZMCreate 生成 6 位数字验证码
func YZMCreate() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06v", rnd.Int31n(1000000))
}

// SleepSecond 秒级别的 sleep
func SleepSecond(duration int) {
	time.Sleep(time.Duration(duration) * time.Second)
}

// SleepMilli 毫秒级别的 sleep
func SleepMilli(duration int) {
	time.Sleep(time.Duration(duration) * time.Millisecond)
}

// ErrorString 解析Error成字符
func ErrorString(err error) string {
	if err != nil {
		errStr := err.Error()
		splist := " desc = "
		if strings.Contains(errStr, splist) {
			splists := strings.Split(errStr, splist)
			if len(splists) == 2 {
				return splists[1]
			}
		}
		return errStr
	}
	return ""
}

// MD5 计算字符串的MD5哈希值
func MD5(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

// MD5Upper 计算字符串的MD5哈希值并返回大写
func MD5Upper(s string) string {
	return strings.ToUpper(MD5(s))
}

// MD5Lower 计算字符串的MD5哈希值并返回小写
func MD5Lower(s string) string {
	return strings.ToLower(MD5(s))
}

// InterfaceToString 将interface{}转换为string
func InterfaceToString(i interface{}) string {
	if i == nil {
		return ""
	}

	switch v := i.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// ToString 简化版本的interface{}转string
func ToString(i interface{}) string {
	return InterfaceToString(i)
}
