## Task: Parse JSON Data with Composed Structs

### Instructions

Write a function `ParseProfile(jsonStr string) (string, int, []string, error)` that takes a JSON string as input and decodes it into a `Profile` struct with a nested `Person` struct (having `Name` string and `Age` int) and a `Hobbies` slice of strings as fields. The function should return the name, age, and hobbies from the JSON data, along with any error that occurs during decoding.

#### Example:
- For a JSON string like `{"person":{"name":"Bob","age":25},"hobbies":["Swimming","Cooking"]}`, the function should return: ("Bob", 25, []string{"Swimming", "Cooking"}, nil)
