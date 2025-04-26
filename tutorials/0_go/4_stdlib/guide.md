# The Go Standard Library

Go’s standard library is a compact toolkit for tasks like JSON data handling, keeping projects lean. Let’s explore `encoding/json` with an example, plus the library’s pros and cons.

## JSON Encoding and Decoding with encoding/json

The `encoding/json` package provides functions for encoding Go structs to JSON and decoding JSON data back into structs.

**Example**:

```go
package main

import (
    "encoding/json"
    "fmt"
)

type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func main() {
    // Create a Person instance
    person := Person{Name: "Alice", Age: 30}
    
    // Encode to JSON
    jsonData, err := json.Marshal(person)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Encoded JSON:", string(jsonData))
    
    // JSON string to decode
    jsonStr := `{"name":"Bob","age":25}`
    
    // Decode JSON to Person struct
    var decodedPerson Person
    err = json.Unmarshal([]byte(jsonStr), &decodedPerson)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Printf("Decoded Person - Name: %s, Age: %d\n", decodedPerson.Name, decodedPerson.Age)
}
```

Run this to see output like:
```
Encoded JSON: {"name":"Alice","age":30}
Decoded Person - Name: Bob, Age: 25
```

## Why Use It?

Go’s standard library is reliable and very complete. You can build full applications with no external dependencies, which will make your code more reliable.
