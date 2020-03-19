package gogtp

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
)

//StreamController 输入输入流控制
type StreamController struct {
	StdIn         io.WriteCloser
	StdOut        io.ReadCloser
	StdErr        io.ReadCloser
	CancelContext context.Context
	CancelFunc    context.CancelFunc
}

//NewStreamController 创建输入输出流控制器
func NewStreamController(in io.WriteCloser, out,err io.ReadCloser) *StreamController {
	return &StreamController{
		StdIn:  in,
		StdOut: out,
		StdErr:err,
	}
}

//Wait 等待命令执行
func (sc *StreamController) Wait() {
	for {
		<-sc.CancelContext.Done()
		break
	}
}

//Stop 停止命令执行
func (sc *StreamController) Stop() {
	if sc.CancelFunc!=nil{
		sc.CancelFunc()
	}
}

//SendCommand 发送命令
func (sc *StreamController) SendCommand(command cmdOptions, respFunc RespFunc) error {
	if sc.CancelContext != nil {
		<-sc.CancelContext.Done()
	}
	in := command.ToString()
	if in == "" {
		return errors.New("command is not null")
	}
	sc.CancelContext, sc.CancelFunc = context.WithCancel(context.TODO())
	firstLine := true
	content := strings.Builder{}
	go func() {
		reader := bufio.NewReader(sc.StdOut)
		for {
			if sc.CancelContext == nil {
				break
			}
			select {
			case <-sc.CancelContext.Done():
				sc.CancelContext = nil
				sc.CancelFunc = nil
				break
			default:
				line, err := reader.ReadString('\n')
				fmt.Println(in,line,err)
				if io.EOF == err {
					break
				} else if err != nil {
					respFunc(Response{
						Command: in,
						Result:  "",
						Error:   err,
					})
					sc.CancelFunc()
				}
				if firstLine && (len(line) == 0 || (!strings.Contains(line, "=") && !strings.Contains(line, "?"))) {
					continue
				}
				firstLine = false
				content.WriteString(line + "\n")
				if command.End {
					respFunc(Response{
						Command: in,
						Result:  content.String(),
						Error:   nil,
					})
					sc.CancelFunc()
					continue
				}
				respFunc(Response{
					Command: in,
					Result:  content.String(),
					Error:   nil,
				})
			}
		}
	}()
	_, err := sc.StdIn.Write([]byte(in+"\n"))
	if err != nil {
		return err
	}
	return nil
}
//ListenStdErr std err 输出监听
func (sc *StreamController) ListenStdErr(sub func(string)) {
	buffer := bufio.NewReader(sc.StdErr)
	go func() {
		for {
			line, err := buffer.ReadString('\n')
			if err != nil || io.EOF == err {
				return
			}
			sub(line)
		}
	}()
}