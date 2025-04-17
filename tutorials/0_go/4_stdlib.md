# The Go Standard Library

Go’s standard library is a compact toolkit for tasks like HTTP servers and JSON handling, keeping projects lean. Let’s explore `encoding/json` and `net/http` with a combined example, plus the library’s pros and cons.

## JSON and HTTP with `encoding/json` and `net/http`

The `encoding/json` package handles JSON marshaling for structs, while `net/http` (improved in Go 1.22) powers web servers.

**Example**:

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func main() {
    // HTTP server with JSON response
    http.HandleFunc("GET /user", func(w http.ResponseWriter, r *http.Request) {
        user := User{Name: "Luna", Age: 25}
        jsonData, _ := json.Marshal(user)
        w.Header().Set("Content-Type", "application/json")
        w.Write(jsonData)
    })

    // Unmarshal JSON for demo
    jsonStr := `{"name":"Max","age":30}`
    var newUser User
    json.Unmarshal([]byte(jsonStr), &newUser)
    fmt.Println("Unmarshaled:", newUser.Name, newUser.Age)

    http.ListenAndServe(":8080", nil)
}
```

Run this, visit `localhost:8080/user` to get `{"name":"Luna","age":25}`, and see “Unmarshaled: Max 30” in the console.

## Why Use It?

Go’s standard library is reliable and very complete. You can build full application with no external dependencies, which will make your code more reliable.
