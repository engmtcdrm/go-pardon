package examples

type Example struct {
	Name string
	Fn   func()
}

var AllExamples = []Example{
	{"Confirm - Basic", ConfirmBasic},
	{"Confirm - Kitchen Sink", ConfirmKitchensink},
	{"Password - Basic", PasswordBasic},
	{"Password - Validate", PasswordValidate},
	{"Password - Kitchen Sink", PasswordKitchesink},
	{"Question - Basic", QuestionBasic},
	{"Question - Validate", QuestionValidate},
	{"Question - Kitchen Sink", QuestionKitchensink},
	{"Select - Basic", SelectBasic},
	{"Select - Struct", SelectStruct},
	{"Select - Kitchen Sink", SelectKitchensink},
}
