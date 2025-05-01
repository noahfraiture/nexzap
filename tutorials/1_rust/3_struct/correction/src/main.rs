mod submission;

#[cfg(test)]
mod tests {
    use super::submission::{Contact, Greet};

    #[test]
    fn test_greet_person() {
        let contact = Contact::Person {
            name: String::from("Alice"),
            phone: String::from("123"),
        };
        assert_eq!(contact.greet(), "Hi, I'm Alice!".to_string());
    }

    #[test]
    fn test_greet_email() {
        let contact = Contact::Email(String::from("bob@example.com"));
        assert_eq!(contact.greet(), "bob@example.com".to_string());  // Uses default
    }
    // Do not test print_greeting
}
