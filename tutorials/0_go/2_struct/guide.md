# Structs and Interfaces in Go

Go skips classes and inheritance, sparing you C++’s maze of virtual destructors or Java’s towering class hierarchies. Instead, it uses **structs** for data, **interfaces** for behavior, and **composition** for reuse, keeping code simple and explicit. Here’s a quick dive with examples.

## Structs: Your Data Blueprint

A **struct** is a custom type grouping fields, like a class without the baggage. Define one with `type` and `struct`, then create instances.

**Example**:

```go
type Person struct {
    Name string
    Age  int
}

p := Person{Name: "Alex", Age: 30}
```

Access fields with dot notation (e.g., `p.Name`). Add methods using a receiver.

**Example**:

```go
func (p Person) Greet() string {
    return "Hi, I’m " + p.Name
}
```

## Composition: Reuse Without Inheritance

Go favors **composition** over Java-style inheritance. Embed a struct to reuse its fields and methods, promoting modularity.

**Example**:

```go
type Employee struct {
    Person
    Role string
}

e := Employee{Person: Person{Name: "Sam", Age: 40}, Role: "Engineer"}
```

`e.Name` and `e.Greet()` work via the embedded `Person`.

## Interfaces: Behavior Contracts

An **interface** lists methods a type must implement, enabling polymorphism. If a type has the methods, it satisfies the interface implicitly.

**Example**:

```go
type Greeter interface {
    Greet() string
}

func SayHello(g Greeter) string {
    return g.Greet()
}
```

`Person` satisfies `Greeter`, so `SayHello(p)` works.

## Why Use Them?

Structs, composition, and interfaces make Go code clean and maintainable, like a well-organized toolbox—simple but effective. The absence of inheritance is a deliberate design choice that reduces unexpected behavior but may require more boilerplate code.
