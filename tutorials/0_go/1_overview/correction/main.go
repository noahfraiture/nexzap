package main

func AddMax(a, b int) int {
	res := a + b
	if res > 255 {
		res = 255
	}
	return res
}
