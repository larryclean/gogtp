package gogtp

import (
	"strings"
)

type cmdOptions struct {
	ID int
	Name string
	Args []string
	End bool
}
//CmdOption 命令设置函数
type CmdOption func(option *cmdOptions)

//CmdID 命令ID
func CmdID(id int) CmdOption {
	return func(option *cmdOptions) {
		option.ID=id
	}
}
//CmdName 命令的名称
func CmdName(name string) CmdOption {
	return func(option *cmdOptions) {
		option.Name=name
	}
}
//CmdArgs 命令的参数
func CmdArgs(args ...string)CmdOption  {
	return func(option *cmdOptions) {
		option.Args=args
	}
}
//CmdEnd 命令是否自动结束
func CmdEnd(end bool)CmdOption  {
	return func(option *cmdOptions) {
		option.End=end
	}
}
//ToString 获取命令执行字符串
func (c *cmdOptions) ToString() string {
	sb := strings.Builder{}
	if c.ID != 0 {
		sb.WriteString(string(c.ID))
	}
	sb.WriteString(c.Name + " ")
	sb.WriteString(strings.Join(c.Args, " "))
	return sb.String()
}
//BuildCommand 命令构建器
func BuildCommand(options ...CmdOption) cmdOptions  {
	op:=cmdOptions{
		End:true,
	}
	for _,v:=range options{
		v(&op)
	}
	return op
}
