package define

import (
	"context"
	"github.com/xwatsonmai/openAgent/agent/aimodel"
	"github.com/xwatsonmai/openAgent/agent/instruction"
)

type IPrompter[I any] interface {
	SystemPrompt(ctx context.Context) (string, error) // 获取系统提示
	StartUserPrompt(ctx context.Context, input I) ([]aimodel.UserContent, error)
	AnswerToInstructions(ctx context.Context, agentAnswer string) ([]instruction.Instruction, error) // 根据Agent的回答解析出指令
}

type IPromptInitializer[I any] interface {
	Initialize(ctx context.Context, input I) (aiChatList aimodel.ChatList, err error) // 初始化Prompt配置
}

type IEntity interface {
	Execute(ctx context.Context, round int, ins instruction.Instruction) ([]aimodel.UserContent, error) // 执行指令
	RoundUserPrompt(ctx context.Context, round int) ([]aimodel.UserContent, error)                      // 获取轮次用户提示
}
