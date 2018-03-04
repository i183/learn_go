package main

import (
	"bufio"
	"fmt"
	"github.com/i183/learn_go/pipeline"
	"os"
)

func main() {
	const filename = "small.in"
	const n = 50
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	ch := pipeline.RandomSource(n)
	writer := bufio.NewWriter(file)
	pipeline.WriterSink(writer, ch)
	writer.Flush()
	file.Close()

	file, err = os.Open(filename)
	if err != nil {
		panic(err)
	}

	ch = pipeline.ReaderSource(bufio.NewReader(file), -1)
	count := 0
	for v := range ch {
		fmt.Println(v)
		count++
		if count > 100 {
			break
		}
	}
	file.Close()

}
