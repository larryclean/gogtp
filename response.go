package gogtp

import (
	"errors"
	"regexp"
	"strings"
)

//Response 响应
type Response struct {
	Command   string
	Result    string
	Error     error
	StdErrOut string
}

//GetResult 获取结果
func (r *Response) GetResult() (string, error) {
	reg, _ := regexp.Compile(`\t`)
	result := strings.TrimSpace(reg.ReplaceAllString(r.Result, " "))

	res := strings.Fields(result)
	l := len(res)
	if l > 0 {
		if res[0] == "=" {
			return strings.TrimSpace(strings.Join(res[1:], "")), nil
		} else if res[0] == "?" {
			return "", errors.New(strings.Join(res[2:], ""))
		}
	}
	return "", errors.New("ERROR: Unrecognized answer: " + result)
}

//RespFunc 输出流回调标准
type RespFunc func(Response)
