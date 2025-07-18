package pardon

type icons struct {
	Alert        string
	QuestionMark string
	Password     string
}

// Default icons available in the package. Overwrite these with
// custom icons using the Icon or IconFunc methods in the respective structs.
var Icons = icons{
	Alert:        "[!] ",
	QuestionMark: "[?] ",
	Password:     "🔒 ",
}
