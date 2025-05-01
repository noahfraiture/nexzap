pub fn process_names(names: Vec<String>) -> Vec<String> {
    names
        .iter()
        .filter(|name| name.len() >= 4)
        .map(|name| name.to_uppercase())
        .collect()
}
