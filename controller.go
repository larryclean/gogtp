package gogtp

import (
	"context"
	"log"
	"os/exec"
	"strings"
)

// EngineOption specifies an option for dialing a Redis server.
type EngineOption func(*engineOptions)

type engineOptions struct {
	name string
	args []string
}

//Controller 控制器
type Controller struct {
	CancelFun        context.CancelFunc
	CancelContext    context.Context
	StreamController *StreamController
	ErrOut           strings.Builder
}

//NewController 创建构建器
func NewController(name string, args ...string) (*Controller, error) {
	ctr := &Controller{
		ErrOut: strings.Builder{},
	}
	ctr.CancelContext, ctr.CancelFun = context.WithCancel(context.TODO())
	cmd := exec.CommandContext(ctr.CancelContext, name, args...)
	stdErr, _ := cmd.StderrPipe()
	stdIn, _ := cmd.StdinPipe()
	stdOut, _ := cmd.StdoutPipe()
	ctr.StreamController = NewStreamController(stdIn, stdOut, stdErr)
	ctr.StreamController.ListenStdErr(func(s string) {
		ctr.ErrOut.WriteString(s)
	})
	err := cmd.Start()
	go func() {
		err := cmd.Wait()
		if err != nil {
			log.Println(err)
		}
		_ = stdErr.Close()
	}()
	return ctr, err
}

//NewControllerByStr 创建构建起根据一个字符串
func NewControllerByStr(path string) (*Controller, error) {
	s1 := strings.Fields(path)
	command := s1[0]
	return NewController(command, s1[1:]...)
}

//Close 关闭
func (c *Controller) Close() {
	c.CancelFun()
}

//StopSendCommand 停止真正之前的程序，开启新的命令
func (c *Controller) StopSendCommand(command cmdOptions, respFunc RespFunc) (err error) {
	c.StreamController.Stop()
	return c.SendCommand(command, respFunc)
}

//SendCommand 发送命令，并使用回调
func (c *Controller) SendCommand(command cmdOptions, respFunc RespFunc) (err error) {
	if err = c.StreamController.SendCommand(command, respFunc); err != nil {
		return
	}
	c.StreamController.Wait()
	return
}

//SyncSendCommand 同步发送命令，并格式化后直接返回结果
func (c *Controller) SyncSendCommand(command cmdOptions) (resp Response, err error) {
	if err = c.StreamController.SendCommand(command, func(response Response) {
		resp = response
	}); err != nil {
		return
	}
	c.StreamController.Wait()
	resp.StdErrOut = c.ResetStdErr()
	return
}

//ResetStdErr 重置std err输出
func (c *Controller) ResetStdErr() string {
	out := c.ErrOut.String()
	c.ErrOut.Reset()
	return out
}
