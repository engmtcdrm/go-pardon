package examples

type Example struct {
	Name string
	Fn   func()
}

var AllExamples = []Example{
	{"1. Confirm - Basic", ConfirmBasic},
	{"2. Password - Basic", PasswordBasic},
	{"3. Question - Basic", QuestionBasic},
	{"4. Select - Basic", SelectBasic},
	{"5. Select - Struct", SelectStruct},
}
