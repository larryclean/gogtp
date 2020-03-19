package gogtp

import (
	"fmt"
	"strings"
)

//Engine GTP 引擎
type Engine struct {
	Controller *Controller
	Name       string
}

//Result 结果
type Result struct {
	Command string
	Result  string
	ErrOut  string
}
//NewEngine 创建简化引擎
func NewEngine(ctr *Controller) (*Engine, error) {
	engine := &Engine{Controller: ctr}
	resp, err := engine.Controller.SyncSendCommand(BuildCommand(CmdName("name")))
	if err != nil {
		return nil, err
	}
	engine.Name, err = resp.GetResult()
	return engine, err
}

//SendCMD 方式自定义命令
func (e *Engine) SendCMD(cmd string) (res Result, err error) {
	cmds := strings.Fields(cmd)
	var resp Response
	if resp, err = e.Controller.SyncSendCommand(BuildCommand(CmdName(cmds[0]), CmdArgs(cmds[1:]...))); err != nil {
		return
	}
	if err = resp.Error; err != nil {
		return
	}
	res.Result, err = resp.GetResult()
	res.ErrOut = resp.StdErrOut
	return
}

//KnowCommand 判断命令是否支持
func (e *Engine) KnowCommand(cmd string) bool {
	value, err := e.SendCMD("known_command " + cmd)
	if err != nil {
		return false
	}
	if strings.ToLower(strings.TrimSpace(value.Result)) != "true" {
		return false
	}
	return true
}

//Komi 设置贴目
func (e *Engine) Komi(komi float64) (Result, error) {
	return e.SendCMD(fmt.Sprintf("komi %v", komi))
}

//BoardSize 设置棋盘大小
func (e *Engine) BoardSize(size int) (Result, error) {
	return e.SendCMD(fmt.Sprintf("boardsize %v", size))
}

//GenMove AI落子
func (e *Engine) GenMove(color string) (Result, error) {
	color = strings.ToUpper(color)
	command := "genmove " + color
	return e.SendCMD(command)
}

//Play 手工落子
func (e *Engine) Play(color, coor string) (Result, error) {
	color = strings.ToUpper(color)
	return e.SendCMD(fmt.Sprintf("play %s %s", color, coor))
}

//LoadSgf 加载SGF文件
func (e *Engine) LoadSgf(file string) (Result, error) {
	command := fmt.Sprintf("loadsgf %s", file)
	return e.SendCMD(command)
}

//FinalStatusList 获取当前盘面形势判断
func (e *Engine) FinalStatusList(cmd string) (Result, error) {
	command := fmt.Sprintf("final_status_list %s", cmd)
	return e.SendCMD(command)
}

//SetLevel 设置AI级别
func (e *Engine) SetLevel(seed int) (Result, error) {
	command := fmt.Sprintf("level %d", seed)
	return e.SendCMD(command)
}

//SetRandomSeed 设置AI随机数
func (e *Engine) SetRandomSeed(seed int) (Result, error) {
	command := fmt.Sprintf("set_random_seed %d", seed)
	return e.SendCMD(command)
}

//ShowBoard 显示棋盘
func (e *Engine) ShowBoard() (Result, error) {
	return e.SendCMD("showboard")
}

//ClearBoard 清空棋盘
func (e *Engine) ClearBoard() (Result, error) {
	return e.SendCMD("clear_board")
}

//PrintSgf 打印SGF
func (e *Engine) PrintSgf() (Result, error) {
	return e.SendCMD("printsgf")
}

//TimeSetting 设置时间规则
func (e *Engine) TimeSetting(baseTime, byoTime, byoStones int) (Result, error) {
	return e.SendCMD(fmt.Sprintf("time_settings %d %d %d", baseTime, byoTime, byoStones))
}

//KGSTimeSetting 设置KGS time
func (e *Engine) KGSTimeSetting(mainTime, readTime, readLimit int) (Result, error) {
	return e.SendCMD(fmt.Sprintf("kgs-time_settings byoyomi %d %d %d", mainTime, readTime, readLimit))
}

//FinalScore 获取结果
func (e *Engine) FinalScore() (Result, error) {
	return e.SendCMD("final_score")
}

//Undo 悔棋
func (e *Engine) Undo() (Result, error) {
	return e.SendCMD("undo")
}

//TimeLeft 设置时间
func (e *Engine) TimeLeft(color string, mainTime, stones int) (Result, error) {
	return e.SendCMD(fmt.Sprintf("time_left %s %d %d", color, mainTime, stones))
}

//Quit 退出
func (e *Engine) Quit() (Result, error) {
	r, err := e.SendCMD("Quit")
	e.Controller.Close()
	return r, err
}
