## Task: Handle Errors in File Reading

In this exercise, you will learn how to handle errors in Go while performing file operations.

### Instructions

Write a function `ReadFileContent(filename string) (string, error)` that takes a filename as input and attempts to read the content of the file using `os.ReadFile`. The function should return the file content as a string if successful, or an empty string and an error if the file cannot be read. Additionally, create a custom error using `errors.New` for a specific case (e.g., if the filename is empty).

#### Example:
- For a valid filename like `data.txt` with content "Hello, Go!", the function should return:
  - Content: "Hello, Go!"
  - Error: nil
- For an empty filename, it should return:
  - Content: ""
  - Error: "filename cannot be empty"
- For a non-existent file, it should return:
  - Content: ""
  - Error: (an error from os.ReadFile indicating file not found)

> **Note**: Do not forget to declare the package with `package main` at the top of the file and import necessary packages.
