package main

import (
	"bufio"
	"fmt"
	"os"
	"pipeline"
	"strconv"
)

func main() {
	//first create pipeline
	p := CreateNetPipeline("large.in", 800000000, 8)
	//写进文件 把 排好序的
	WriteFile(p, "large.out")
	//打印文件中已经排好序的
	PrintResult("large.out")
}
func PrintResult(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	p := pipeline.ReaderSource(bufio.NewReader(file), -1)
	count := 0
	for v := range p {
		count++
		fmt.Println(v)
		if count > 20 {
			break
		}

	}
}
func WriteFile(p <-chan int, fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	//defer 是先进后出的
	defer file.Close()
	writer := bufio.NewWriter(file)
	pipeline.WriteSink(writer, p)
	defer writer.Flush()
}
func CreatePipeline(fileName string, fileSize int, chunkCount int) <-chan int {
	chunkSize := fileSize / chunkCount
	sorts := []<-chan int{}
	pipeline.Init()

	for i := 0; i < chunkCount; i++ {
		file, err := os.Open(fileName)
		if err != nil {
			panic(err)
		}
		file.Seek(int64(i*chunkSize), 0)
		source := pipeline.ReaderSource(bufio.NewReader(file), chunkSize)
		sorts = append(sorts,  pipeline.InMeSort(source))
	}
	return pipeline.MergeN(sorts...)
}
func CreateNetPipeline(fileName string, fileSize int, chunkCount int) <-chan int {
	chunkSize := fileSize / chunkCount
	sorts := []<-chan int{}
	addrStr := []string{}
	pipeline.Init()

	for i := 0; i < chunkCount; i++ {
		file, err := os.Open(fileName)
		if err != nil {
			panic(err)
		}
		file.Seek(int64(i*chunkSize), 0)
		source := pipeline.ReaderSource(bufio.NewReader(file), chunkSize)
		addr:=":"+strconv.Itoa(7000+i)
		pipeline.NetWorkSink(addr, pipeline.InMeSort(source))
		addrStr=append(addrStr,addr)
	}
	for _,addr:=range addrStr  {
		sorts = append(sorts,  pipeline.NetWorkResource(addr))
	}
	return pipeline.MergeN(sorts...)
}