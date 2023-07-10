

class ChatGeneration(Generation):
"""Output of a single generation."""

    text = ""
    message: BaseMessage

    @root_validator
    def set_text(cls, values: Dict[str, Any]) -> Dict[str, Any]:
        values["text"] = values["message"].content
        return values

class ChatResult(BaseModel):
"""Class that contains all relevant information for a Chat Result."""
# just a list of ai responses
generations: List[ChatGeneration]
"""List of the things generated."""
llm_output: Optional[dict] = None
"""For arbitrary LLM provider specific output."""
