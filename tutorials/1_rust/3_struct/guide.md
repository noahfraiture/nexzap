# Structs, Enums, Traits, and Impl in Rust

Rust’s data and behavior tools—**structs**, **enums**, **traits**, and **impl**—are like a compact LEGO set: simple pieces for building flexible, type-safe code without classes. Here’s a quick, beginner-friendly overview with an example.

## The Core Pieces

- **Structs** bundle related data, like a character sheet.
- **Enums** define types with variants, great for choices. Variants can contain values.
- **Traits** list methods for shared behaviors, like a contract. It can have a default behavior
- **Impl** attaches methods or trait implementations to structs/enums.

**Example**:

```rust
struct Person {
    name: String,
    age: u32,
}

enum Message {
    Text(String),
    Image,
}

trait Greet {
    fn greet(&self) -> String {
        String::from("Default implementation")
    }
}

impl Greet for Person {
}

impl Greet for Message {
    fn greet(&self) -> String {
        match self {
            Message::Text(t) => t.clone(),
            Message::Image => String::from("This is an image"),
        }
    }
}

fn main() {
    let p = Person { name: String::from("Sam"), age: 25 };
    let m = Message::Text(String::from("Hi!"));
    println!("{}", p.greet()); // Default implementation
    println!("{}", m.greet()); // Hi!
}
```

## Why Use Them?

These tools make Rust modular and safe, perfect for CLI tools or game engines. Downside? They require more code than class-based languages. Still, they’re a powerful combo for reliable code.

Snap them together and build something cool!
