package pipeline

import (
	"encoding/binary"
	"io"
	"math/rand"
	"sort"
)

func ArraySource(numbers ...int) <-chan int {
	out := make(chan int, 1024)
	go func() {
		for _, v := range numbers {
			out <- v
		}
		close(out)
	}()
	return out
}

func InMemSort(in <-chan int) <-chan int {
	out := make(chan int, 1024)
	go func() {
		var numbers []int
		for v := range in {
			numbers = append(numbers, v)
		}

		sort.Ints(numbers)

		count := 0
		for _, v := range numbers {
			out <- v
			count++
		}
		close(out)
	}()
	return out
}

func Merge(in1, in2 <-chan int) <-chan int {
	out := make(chan int, 1024)
	go func() {
		v1, ok1 := <-in1
		v2, ok2 := <-in2
		for ok1 || ok2 {
			if !ok2 || (ok1 && v1 <= v2) {
				out <- v1
				v1, ok1 = <-in1
			} else {
				out <- v2
				v2, ok2 = <-in2
			}
		}
		close(out)
	}()
	return out
}

func ReaderSource(reader io.Reader, chunkSize int) <-chan int {
	out := make(chan int, 1024)
	go func() {
		buffer := make([]byte, 8)
		bytesRead := 0
		for {
			n, err := reader.Read(buffer)
			bytesRead += n
			if n > 0 {
				v := int(binary.BigEndian.Uint64(buffer))
				out <- v
			}

			if err != nil || (chunkSize != -1 && bytesRead >= chunkSize) {
				break
			}
		}
		close(out)
	}()
	return out
}

func WriterSink(writer io.Writer, in <-chan int) {
	buffer := make([]byte, 8)
	count := 0
	for v := range in {
		binary.BigEndian.PutUint64(buffer, uint64(v))
		writer.Write(buffer)
		count++
	}
}

func RandomSource(count int) <-chan int {
	out := make(chan int, 1024)
	go func() {
		for i := 0; i < count; i++ {
			out <- rand.Int()
		}
		close(out)
	}()
	return out
}

func MergeN(inputs ...<-chan int) <-chan int {
	if len(inputs) == 1 {
		return inputs[0]
	}

	m := len(inputs) / 2

	return Merge(MergeN(inputs[:m]...), MergeN(inputs[m:]...))
}
