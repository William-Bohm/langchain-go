package qa

import (
	"fmt"
	"promptSchema"
	"strings"
)

func main() {
	template := `You are a teacher grading a quiz.
You are given a question, the student's answer, and the true answer, and are asked to score the student answer as either CORRECT or INCORRECT.

Example Format:
QUESTION: question here
STUDENT ANSWER: student's answer here
TRUE ANSWER: true answer here
GRADE: CORRECT or INCORRECT here

Grade the student answers based ONLY on their factual accuracy. Ignore differences in punctuation and phrasing between the student answer and true answer. It is OK if the student answer contains more information than the true answer, as long as it does not contain any conflicting statements. Begin! 

QUESTION: %s
STUDENT ANSWER: %s
TRUE ANSWER: %s
GRADE:`
	prompt := promptSchema.NewPromptTemplate([]string{"query", "result", "answer"}, template)

	contextTemplate := `You are a teacher grading a quiz.
You are given a question, the context the question is about, and the student's answer. You are asked to score the student's answer as either CORRECT or INCORRECT, based on the context.

Example Format:
QUESTION: question here
CONTEXT: context the question is about here
STUDENT ANSWER: student's answer here
GRADE: CORRECT or INCORRECT here

Grade the student answers based ONLY on their factual accuracy. Ignore differences in punctuation and phrasing between the student answer and true answer. It is OK if the student answer contains more information than the true answer, as long as it does not contain any conflicting statements. Begin! 

QUESTION: %s
CONTEXT: %s
STUDENT ANSWER: %s
GRADE:`
	contextPrompt := NewPromptTemplate([]string{"query", "context", "result"}, contextTemplate)

	cotTemplate := `You are a teacher grading a quiz.
You are given a question, the context the question is about, and the student's answer. You are asked to score the student's answer as either CORRECT or INCORRECT, based on the context.
Write out in a step by step manner your reasoning to be sure that your conclusion is correct. Avoid simply stating the correct answer at the outset.

Example Format:
QUESTION: question here
CONTEXT: context the question is about here
STUDENT ANSWER: student's answer here
EXPLANATION: step by step reasoning here
GRADE: CORRECT or INCORRECT here

Grade the student answers based ONLY on their factual accuracy. Ignore differences in punctuation and phrasing between the student answer and true answer. It is OK if the student answer contains more information than the true answer, as long as it does not contain any conflicting statements. Begin! 

QUESTION: %s
CONTEXT: %s
STUDENT ANSWER: %s
EXPLANATION:`
	cotPrompt := NewPromptTemplate([]string{"query", "context", "result"}, cotTemplate)

	// Print the templates to make sure they are created correctly.
	fmt.Println(strings.ReplaceAll(prompt.Template, "%s", "{}"))
	fmt.Println(strings.ReplaceAll(contextPrompt.Template, "%s", "{}"))
	fmt.Println(strings.ReplaceAll(cotPrompt.Template, "%s", "{}"))
}
