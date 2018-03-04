package main

import (
	"bufio"
	"fmt"
	"github.com/i183/learn_go/pipeline"
	"os"
	"time"
)

func main() {
	const inFilename = "large.in"
	const outFilename = "large.out"
	const chunkCount = 10

	startTime := time.Now()

	ch := createPipeline(inFilename, chunkCount)
	writeToFIle(ch, outFilename)
	//printFile(outFilename)

	fmt.Println(time.Now().Sub(startTime))
}

func printFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	ch := pipeline.ReaderSource(file, -1)
	count := 0
	for v := range ch {
		fmt.Println(v)

		count++
		if count >= 100 {
			break
		}
	}

}

func writeToFIle(ch <-chan int, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	pipeline.WriterSink(writer, ch)
	writer.Flush()
}

func createPipeline(filename string, chunkCount int) <-chan int {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		panic(err)
	}
	sortResults := []<-chan int{}
	chunkSizes := segmentation(int(fileInfo.Size()), chunkCount)
	fmt.Println(chunkSizes)
	readLocation := 0
	for i := 0; i < chunkCount; i++ {
		file, err := os.Open(filename)
		if err != nil {
			panic(err)
		}

		chunkSize := chunkSizes[i] * 8
		file.Seek(int64(readLocation), 0)
		readLocation += chunkSize

		source := pipeline.ReaderSource(bufio.NewReader(file), chunkSize)
		sortResults = append(sortResults, pipeline.InMemSort(source))
	}

	return pipeline.MergeN(sortResults...)
}

/**
分配每块处理Int数
*/
func segmentation(fileSize, chunkCount int) []int {
	intCountOfFile := fileSize / 8
	chunkSize, over := intCountOfFile/chunkCount, intCountOfFile%chunkCount
	chunkSizes := make([]int, chunkCount)
	for i := 0; i < chunkCount; i++ {
		if i < over {
			chunkSizes[i] = chunkSize + 1
		} else {
			chunkSizes[i] = chunkSize
		}
	}
	return chunkSizes[:]
}
