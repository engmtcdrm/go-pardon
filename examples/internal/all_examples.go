package internal

import "github.com/engmtcdrm/go-eggy"

var AllExamples = []eggy.Example{
	{Name: "New Question - Basic", Fn: ExampleNewQuestion},
	{Name: "New Question - Validate", Fn: NewQuestionValidate},
	{Name: "New Question - Kitchen Sink", Fn: NewQuestionKitchensink},
	{Name: "Confirm - Basic", Fn: ConfirmBasic},
	{Name: "Confirm - Kitchen Sink", Fn: ConfirmKitchensink},
	{Name: "Password - Basic", Fn: PasswordBasic},
	{Name: "Password - Validate", Fn: PasswordValidate},
	{Name: "Password - Kitchen Sink", Fn: PasswordKitchesink},
	{Name: "Question - Basic", Fn: QuestionBasic},
	{Name: "Question - Validate", Fn: QuestionValidate},
	{Name: "Question - Kitchen Sink", Fn: QuestionKitchensink},
	{Name: "Select - Basic", Fn: SelectBasic},
	{Name: "Select - Struct", Fn: SelectStruct},
	{Name: "Select - Kitchen Sink", Fn: SelectKitchensink},
	{Name: "Form - Basic", Fn: FormBasic},
	{Name: "Form - Validate", Fn: FormValidate},
	{Name: "Reusing a Prompt", Fn: ReusingPrompt},
}
