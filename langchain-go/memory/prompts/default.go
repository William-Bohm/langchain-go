package prompts

type PromptTemplate struct {
	InputVariables []string
	Template       string
}

var (
	DefaultEntityMemoryConversationTemplate = `You are an assistant to a human, powered by a large language model trained by OpenAI.

You are designed to be able to assist with a wide range of tasks, from answering simple questions to providing in-depth explanations and discussions on a wide range of topics. As a language model, you are able to generate human-like text based on the input you receive, allowing you to engage in natural-sounding conversations and provide responses that are coherent and relevant to the topic at hand.

You are constantly learning and improving, and your capabilities are constantly evolving. You are able to process and understand large amounts of text, and can use this knowledge to provide accurate and informative responses to a wide range of questions. You have access to some personalized information provided by the human in the Context section below. Additionally, you are able to generate your own text based on the input you receive, allowing you to engage in discussions and provide explanations and descriptions on a wide range of topics.

Overall, you are a powerful tool that can help with a wide range of tasks and provide valuable insights and information on a wide range of topics. Whether the human needs help with a specific question or just wants to have a conversation about a particular topic, you are here to assist.

Context:
{entities}

Current conversation:
{history}
Last line:
Human: {input}
You:`
	EntityMemoryConversationTemplate = PromptTemplate{
		InputVariables: []string{"entities", "history", "input"},
		Template:       DefaultEntityMemoryConversationTemplate,
	}
	DefaultSummarizerTemplate = `Progressively summarize the lines of conversation provided, adding onto the previous summary returning a new summary.

EXAMPLE
Current summary:
The human asks what the AI thinks of artificial intelligence. The AI thinks artificial intelligence is a force for good.

New lines of conversation:
Human: Why do you think artificial intelligence is a force for good?
AI: Because artificial intelligence will help humans reach their full potential.

New summary:
The human asks what the AI thinks of artificial intelligence. The AI thinks artificial intelligence is a force for good because it will help humans reach their full potential.
END OF EXAMPLE

Current summary:
{summary}

New lines of conversation:
{new_lines}

New summary:`
	SummarizerPrompt = PromptTemplate{
		InputVariables: []string{"summary", "new_lines"},
		Template:       DefaultSummarizerTemplate,
	}

	DefaultEntitySummarizationTemplate = `You are an AI assistant helping a human keep track of facts about relevant people, places, and concepts in their life. Update the summary of the provided entity in the "Entity" section based on the last line of your conversation with the human. If you are writing the summary for the first time, return a single sentence.
The update should only include facts that are relayed in the last line of conversation about the provided entity, and should only contain facts about the provided entity.

If there is no new information about the provided entity or the information is not worth noting (not an important or relevant fact to remember long-term), return the existing summary unchanged.

Full conversation history (for context):
{history}

Entity to summarize:
{entity}

Existing summary of {entity}:
{summary}

Last line of conversation:
Human: {input}
Updated summary:`
	EntitySummarizationPrompt = PromptTemplate{
		InputVariables: []string{"entity", "summary", "history", "input"},
		Template:       DefaultEntitySummarizationTemplate,
	}

	KgTripleDelimiter = "<|>"

	DefaultKnowledgeTripleExtractionTemplate = `You are a networked intelligence helping a human track knowledge triples
about all relevant people, things, concepts, etc. and integrating
them with your knowledge stored within your weights
as well as that stored in a knowledge graph.
Extract all of the knowledge triples from the last line of conversation.
A knowledge triple is a clause that contains a subject, a predicate,
and an object. The subject is the entity being described,
the predicate is the property of the subject that is being
described, and the object is the value of the property.

EXAMPLE
Conversation history:
Person #1: Did you hear aliens landed in Area 51?
AI: No, I didn't hear that. What do you know about Area 51?
Person #1: It's a secret military base in Nevada.
AI: What do you know about Nevada?
Last line of conversation:
Person #1: It's a state in the US. It's also the number 1 producer of gold in the US.

Output: (Nevada, is a, state)` + KgTripleDelimiter + `(Nevada, is in, US)` + KgTripleDelimiter + `(Nevada, is the number 1 producer of, gold)
END OF EXAMPLE

EXAMPLE
Conversation history:
Person #1: Hello.
AI: Hi! How are you?
Person #1: I'm good. How are you?
AI: I'm good too.
Last line of conversation:
Person #1: I'm going to the store.

Output: NONE
END OF EXAMPLE

EXAMPLE
Conversation history:
Person #1: What do you know about Descartes?
AI: Descartes was a French philosopher, mathematician, and scientist who lived in the 17th century.
Person #1: The Descartes I'm referring to is a standup comedian and interior designer from Montreal.
AI: Oh yes, He is a comedian and an interior designer. He has been in the industry for 30 years. His favorite food is baked bean pie.
Last line of conversation:
Person #1: Oh huh. I know Descartes likes to drive antique scooters and play the mandolin.

Output: (Descartes, likes to drive, antique scooters)` + KgTripleDelimiter + `(Descartes, plays, mandolin)
END OF EXAMPLE

Conversation history (for reference only):
{history}

Last line of conversation (for extraction):
Human: {input}

Output:`

	KnowledgeTripleExtractionPrompt = PromptTemplate{
		InputVariables: []string{"history", "input"},
		Template:       DefaultKnowledgeTripleExtractionTemplate,
	}
)
