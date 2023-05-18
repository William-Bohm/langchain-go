package toolSchema

type InvalidTool struct {
	BaseTool
	Name        string
	Description string
}

func NewInvalidTool() *InvalidTool {
	base := &BaseTool{Name: "invalid_tool"}
	return &InvalidTool{BaseTool: *base, Description: "Called when tool name is invalid."}
}

func (it *InvalidTool) Run(toolName string) string {
	return toolName + " is not a valid tool, try another one."
}

func (it *InvalidTool) ARun(toolName string) string {
	return toolName + " is not a valid tool, try another one."
}
