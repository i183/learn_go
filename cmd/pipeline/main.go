package main

import (
	"fmt"
	"github.com/i183/learn_go/pipeline"
)

func main() {
	ch := pipeline.ArraySource(8, 3, 5, 3, 1, 9)

	for v := range ch {
		fmt.Println(v)
	}
}
