package document_compressor

// promptTemplate is a constant string used as a template.
const chainFilterPromptTemplate = `Given the following question and context, return YES if the context is relevant to the question and NO if it isn't.

> Question: {question}
> Context:
>>>
{context}
>>>
> Relevant (YES / NO):`
