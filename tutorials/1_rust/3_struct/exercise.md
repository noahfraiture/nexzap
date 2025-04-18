## Task: Implement a Greetable Contact System

In this exercise, you will use structs, enums, traits, and impl to model a contact system where different entities can greet using a trait’s default and custom implementations.

### Instructions

Define a `Contact` enum with variants `Person { name: String, phone: String }` and `Email(String)`. Create a `Greet` trait with a `greet(&self) -> String` method that has a default implementation returning `"Hello!"`. Implement `Greet` for `Contact` such that `Person` returns a greeting with their name (e.g., `"Hi, I'm Alice!"`) and `Email` uses the default. Write a function `print_greeting(c: &Contact)` that prints the greeting.

#### Example:
- Input: `Contact::Person { name: String::from("Alice"), phone: String::from("123") }`, `Contact::Email(String::from("bob@example.com"))`
- Output: `Hi, I'm Alice!`, `Hello!`

> **Note**: Ensure `fn main` is defined and test both variants.
