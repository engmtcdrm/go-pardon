package examples

type Example struct {
	Name string
	Fn   func()
}

var AllExamples = []Example{
	{"1. Select - Basic", SelectBasic},
	{"2. Select - Struct", SelectStruct},
	{"3. Confirm - Basic", ConfirmBasic},
	{"4. Question - Basic", QuestionBasic},
	{"5. Password - Basic", PasswordBasic},
}
