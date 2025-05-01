pub fn string_length(input: &str) -> usize {
    let mut res = 0;
    for _c in input.chars() {
        res += 1;
    }
    res
}
