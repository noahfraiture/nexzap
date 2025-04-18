package main

func CalculateSum(numbers []int) int {
	if len(numbers) == 0 {
		return 0
	}
	if len(numbers) == 1 {
		return numbers[0]
	}

	ch := make(chan int)
	mid := len(numbers) / 2

	go func() {
		sum := 0
		for _, num := range numbers[:mid] {
			sum += num
		}
		ch <- sum
	}()

	go func() {
		sum := 0
		for _, num := range numbers[mid:] {
			sum += num
		}
		ch <- sum
	}()

	return <-ch + <-ch
}
