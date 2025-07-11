package instruction

type Type string

func (i Type) ToString() string {
	switch i {
	case TypeEnd:
		return "结束指令"
	}
	return string(i)
}

const (
	TypeEnd = "end" // 结束指令
)

type Instruction struct {
	ID      string `json:"id"`
	Type    Type   `json:"type"`    // 指令类型: message, takeover, mcp_tool
	Target  string `json:"target"`  // 指令目标: 如mcp工具调用目标
	Content any    `json:"content"` // 指令内容: 消息内容、接管信息或MCP工具的入参
}

func (i Instruction) IsEnd() bool {
	return i.Type == TypeEnd
}
