package callbackSchema

import (
	"github.com/William-Bohm/langchain-go/langchain-go/llm/llmSchema"
	"github.com/google/uuid"
	"time"
)

type TracerSessionBase struct {
	StartTime time.Time
	Name      *string
	Extra     *map[string]interface{}
}

type TracerSessionCreate struct {
	TracerSessionBase
}

type TracerSession struct {
	TracerSessionBase
	Id int
}

type BaseRun struct {
	Id             *interface{}
	StartTime      time.Time
	EndTime        time.Time
	Extra          *map[string]interface{}
	ExecutionOrder int
	Serialized     map[string]interface{}
	SessionId      uuid.UUID
	Error          *error
}

func NewBaseRun(executionOrder int, serialized map[string]interface{}, sessionId uuid.UUID, extra *map[string]interface{}) BaseRun {
	return BaseRun{
		StartTime:      time.Now(),
		EndTime:        time.Now(),
		ExecutionOrder: executionOrder,
		Serialized:     serialized,
		SessionId:      sessionId,
		Extra:          extra,
	}
}

type LLMRun struct {
	BaseRun
	Prompts  []string
	Response *llmSchema.LLMResult
}

func NewLLMRun(executionOrder int, serialized map[string]interface{}, sessionId uuid.UUID, prompts []string, extra *map[string]interface{}) LLMRun {
	baseRun := NewBaseRun(executionOrder, serialized, sessionId, extra)
	return LLMRun{
		BaseRun: baseRun,
		Prompts: prompts,
	}
}

type ChainRun struct {
	BaseRun
	Inputs         map[string]interface{}
	Outputs        *map[string]interface{}
	ChildLLMRuns   []LLMRun
	ChildChainRuns []ChainRun
	ChildToolRuns  []ToolRun
	ChildRuns      []interface{} // Can be LLMRun, ChainRun, or ToolRun
}

func NewChainRun(executionOrder int, serialized map[string]interface{}, sessionId uuid.UUID, inputs map[string]interface{}, extra *map[string]interface{}) ChainRun {
	baseRun := NewBaseRun(executionOrder, serialized, sessionId, extra)
	return ChainRun{
		BaseRun: baseRun,
		Inputs:  inputs,
	}
}

type ToolRun struct {
	BaseRun
	ToolInput      string
	Output         *string
	Action         string
	ChildLLMRuns   []LLMRun
	ChildChainRuns []ChainRun
	ChildToolRuns  []ToolRun
	ChildRuns      []interface{} // Can be LLMRun, ChainRun, or ToolRun
}

// variable 'action' is the string version of serialized
func NewToolRun(executionOrder int, serialized map[string]interface{}, sessionId uuid.UUID, toolInput string, extra *map[string]interface{}) ToolRun {
	baseRun := NewBaseRun(executionOrder, serialized, sessionId, extra)
	return ToolRun{
		BaseRun:   baseRun,
		ToolInput: toolInput,
	}
}
