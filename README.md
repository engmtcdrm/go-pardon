
# go-pardon

Lightweight interactive CLI prompt library for Go

## Import the package

```go
import "github.com/engmtcdrm/pardon"
```

## Usage

### Confirm Prompt

```go
continueFlag := true
confirm := pardon.NewConfirm(&continueFlag).
    Title("Are you sure you want to proceed?")

if err := confirm.Ask(); err != nil {
    fmt.Printf("Error: %v\n", err)
}
if continueFlag {
    fmt.Println("Proceeding!")
} else {
    fmt.Println("Stopping!")
}
```

### Password Prompt

```go
password := []byte{}
passwordPrompt := pardon.NewPassword(&password).
    Title("Enter your password:")

if err := passwordPrompt.Ask(); err != nil {
    fmt.Printf("Error: %v\n", err)
}
fmt.Printf("Entered password is '%s'\n", string(password))
```

### Question Prompt

```go
favColor := ""
question := pardon.NewQuestion(&favColor).
    Title("What is your favorite color?")

if err := question.Ask(); err != nil {
    fmt.Printf("Error: %v\n", err)
}
fmt.Printf("Entered favorite color is '%s'\n", favColor)
```

### Select Prompt

```go
package main

import (
    "fmt"
    "github.com/engmtcdrm/pardon"
)

func main() {
    var selectedColor int
    colors := []pardon.Option[int]{
        {Key: "Red", Value: 1},
        {Key: "Blue", Value: 2},
        {Key: "Green", Value: 3},
        {Key: "Yellow", Value: 4},
    }

    selectPrompt := pardon.NewSelect(&selectedColor).
        Title("Choose a color:").
        Options(colors...)

    if err := selectPrompt.Ask(); err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("Selected option: %v\n", selectedColor)
}
```

### Form (Multiple Prompts)

```go
package main

import (
    "fmt"
    "github.com/engmtcdrm/pardon"
)

func main() {
    // Variables to store results
    var name string
    var password []byte
    var age int
    var confirmed bool

    // Create age options
    ageOptions := []pardon.Option[int]{
        {Key: "18-25", Value: 22},
        {Key: "26-35", Value: 30},
        {Key: "36-45", Value: 40},
        {Key: "46+", Value: 50},
    }

    // Create form with multiple prompts
    form := pardon.NewForm(
        pardon.NewQuestion(&name).
            Title("What is your name?"),

        pardon.NewPassword(&password).
            Title("Enter your password:"),

        pardon.NewSelect(&age).
            Title("Select your age group:").
            Options(ageOptions...),

        pardon.NewConfirm(&confirmed).
            Title("Do you agree to the terms?"),
    )

    // Execute all prompts in sequence
    if err := form.Ask(); err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    // Display results
    fmt.Printf("Name: %s\n", name)
    fmt.Printf("Password length: %d characters\n", len(password))
    fmt.Printf("Age group: %d\n", age)
    fmt.Printf("Agreed to terms: %t\n", confirmed)
}
```
