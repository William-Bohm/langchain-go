package document_compressor

// promptTemplate is a constant string used as a template.
const chainExtractPromptTemplate = `Given the following question and context, extract any part of the context *AS IS* that is relevant to answer the question. If none of the context is relevant return {no_output_str}. 

Remember, *DO NOT* edit the extracted parts of the context.

> Question: {{question}}
> Context:
>>>
{{context}}
>>>
Extracted relevant parts:`

// noOutputStr is a constant string used to represent the output when no relevant parts are found.
const noOutputStr = "{no_output_str}"
