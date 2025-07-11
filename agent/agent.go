package agent

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/xwatsonmai/openAgent/agent/agentError"
	"github.com/xwatsonmai/openAgent/agent/aimodel"
	"github.com/xwatsonmai/openAgent/agent/define"
)

type Agent[I any] struct {
	PromptConfig json.RawMessage     `json:"prompt_config"` // Prompt配置
	prompter     define.IPrompter[I] // Prompter接口
	entity       define.IEntity      // 实体接口
	ai           aimodel.IAiModel
	round        int // 当前轮次
	aiChatList   aimodel.ChatList
}

func New[I any](prompter define.IPrompter[I], entity define.IEntity, ai aimodel.IAiModel) *Agent[I] {
	return &Agent[I]{
		prompter:   prompter,
		entity:     entity,
		ai:         ai,
		round:      0,
		aiChatList: aimodel.ChatList{},
	}
}

func (a *Agent[I]) Do(ctx context.Context, input I) error {
	// 检查看prompter是否实现了IPromptInitializer接口，如果实现了，说明上层业务需要自行初始化与Agent的对话消息
	if initializer, ok := a.prompter.(define.IPromptInitializer[I]); ok {
		// 框架只依赖IPromptInitializer接口，具体的实现由上层业务提供
		aiChatList, err := initializer.Initialize(ctx, input)
		if err != nil {
			return err
		}
		a.aiChatList = aiChatList
	} else {
		// 没有实现，则使用默认的初始化方式
		systemPrompt, err := a.prompter.SystemPrompt(ctx)
		if err != nil {
			return err
		}
		startUserPrompt, err := a.prompter.StartUserPrompt(ctx, input)
		if err != nil {
			return err
		}
		startRoundUserPrompt, err := a.entity.RoundUserPrompt(ctx, a.round)
		if err != nil {
			return err
		}
		// 把startUserPrompt和startRoundUserPrompt合并
		allUserPrompt := append(startUserPrompt, startRoundUserPrompt...)
		a.aiChatList = aimodel.ChatList{
			{
				Role:    aimodel.EAIChatRoleSystem,
				Content: systemPrompt,
			},
			{
				Role:    aimodel.EAIChatRoleUser,
				Content: allUserPrompt,
			},
		}
	}

	for {
		a.round++
		// 暂时只支持非流式对话
		aiResult, err := a.ai.Chat(ctx, a.aiChatList)
		if err != nil {
			return err
		}
		answer := aiResult.Result
		instructions, err := a.prompter.AnswerToInstructions(ctx, answer)
		if err != nil {
			return err
		}
		if len(instructions) == 0 {
			// 如果没有指令，可能是解析异常了
			return errors.New("no instructions found in agent answer")
		}
		var thisRoundUserPrompt []aimodel.UserContent
		// 执行指令
		for _, ins := range instructions {
			if ins.IsEnd() {
				return agentError.End
			}
			// 执行指令
			userContent, err := a.entity.Execute(ctx, a.round, ins)
			if err != nil {
				return err
			}
			thisRoundUserPrompt = append(thisRoundUserPrompt, userContent...)

		}
		roundUserPrompt, err := a.entity.RoundUserPrompt(ctx, a.round)
		if err != nil && errors.Is(err, agentError.End) {
			return agentError.End
		}
		if err != nil {
			return err
		}
		// 把thisRoundUserPrompt和roundUserPrompt合并
		allUserPrompt := append(thisRoundUserPrompt, roundUserPrompt...)
		a.aiChatList = append(a.aiChatList, aimodel.Chat{
			Role:    aimodel.EAIChatRoleUser,
			Content: allUserPrompt,
		})
	}
}

//func (a *Agent[I]) Result() R {
//	return a.result
//}

//func (a *Agent[I, R]) aiChat(ctx context.Context, aiChatList aimodel.ChatList, flow bool) {}
