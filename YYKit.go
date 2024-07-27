package YY

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

var DEBUG bool = true

func Str(i int) string {
	return string(i)
}
func StrInt32(i int32) string {
	return string(i)
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

func CreateNumber(bit int) string {
	var result = ""
	for i := 0; i < bit; i++ {
		result = result + numberString[rand.Int31n(10)]
	}
	return result
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
