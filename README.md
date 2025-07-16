
# go-pardon
Lightweight interactive CLI prompt library for Go

![](https://media.giphy.com/media/Nmc3muJhaCfPe2LWd9/giphy.gif)

## Import the package
```go
import "github.com/engmtcdrm/pardon"
```

## Usage

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

    selectPrompt := pardon.NewSelect[int]().
        Title("Choose a color:").
        Options(colors...).
        Value(&selectedColor)

    if err := selectPrompt.Ask(); err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("Selected option: %v\n", selectedColor)
}
```

### Question Prompt
```go
favColor := ""
question := pardon.NewQuestion().
    Title("What is your favorite color?").
    Value(&favColor)

if err := question.Ask(); err != nil {
    fmt.Printf("Error: %v\n", err)
}
fmt.Printf("Entered favorite color is '%s'\n", favColor)
```

### Password Prompt
```go
password := []byte{}
passwordPrompt := pardon.NewPassword().
    Title("Enter your password:").
    Value(&password)

if err := passwordPrompt.Ask(); err != nil {
    fmt.Printf("Error: %v\n", err)
}
fmt.Printf("Entered password is '%s'\n", string(password))
```

### Confirm Prompt
```go
continueFlag := true
confirm := pardon.NewConfirm().
    Title("Are you sure you want to proceed?").
    Value(&continueFlag)

if err := confirm.Ask(); err != nil {
    fmt.Printf("Error: %v\n", err)
}
if continueFlag {
    fmt.Println("Proceeding!")
} else {
    fmt.Println("Stopping!")
}
```
