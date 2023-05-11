implement all of the functions for me, dont include the whole file
Here is some reference code (all structs have a New____Run() function that creates on for you):

type BaseRun struct {
Id             *interface{}
StartTime      time.Time
EndTime        time.Time
Extra          *Any
ExecutionOrder int
Serialized     Any
SessionId      int
Error          *string
}

type LLMRun struct {
BaseRun
Prompts  []string
Response *LLMResult
}

type ChainRun struct {
BaseRun
Inputs         Any
Outputs        *Any
ChildLLMRuns   []LLMRun
ChildChainRuns []ChainRun
ChildToolRuns  []ToolRun
ChildRuns      []interface{} // Can be LLMRun, ChainRun, or ToolRun
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