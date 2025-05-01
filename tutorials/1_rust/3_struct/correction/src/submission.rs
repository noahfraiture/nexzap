pub enum Contact {
    Person { name: String, phone: String },
    Email(String),
}

pub trait Greet {
    fn greet(&self) -> String {
        "Hello!".to_string()
    }
}

impl Greet for Contact {
    fn greet(&self) -> String {
        match self {
            Contact::Person { name, .. } => format!("Hi, I'm {}!", name),
            Contact::Email(s) => s.clone(),  // Use the default implementation
        }
    }
}
