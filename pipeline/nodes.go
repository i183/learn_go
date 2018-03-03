package pipeline

func ArraySource(numbers ...int) chan int {
	out := make(chan int)
	go func() {
		for _, v := range numbers {
			out <- v
		}
		close(out)
	}()
	return out
}
