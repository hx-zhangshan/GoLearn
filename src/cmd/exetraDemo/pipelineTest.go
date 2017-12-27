package main

import (
	"bufio"
	"os"
	"pipeline"
	"fmt"
)

func main() {
	//first create pipeline
	p := CreatePipeline("small.txt", 512, 4)
	//写进文件 把 排好序的
	WriteFile(p,"smallout.txt")
	//打印文件中已经排好序的
	PrintResult("smallout.txt")
}
func PrintResult(fileName string) {
	file,err:=os.Open(fileName)
	if err!=nil{
		panic(err)
	}
	defer file.Close()
	p := pipeline.ReaderSource(bufio.NewReader(file), -1)
	for v:=range p  {
		fmt.Println(v)
	}
}
func WriteFile(p <-chan int,fileName string) {
	file,err:=os.Create(fileName)
	if err!=nil {
		panic(err)
	}
	//defer 是先进后出的
	defer file.Close()
	writer := bufio.NewWriter(file)
	pipeline.WriteSink(writer,p)
	defer writer.Flush()
}
func CreatePipeline(fileName string, fileSize int, chunkCount int) <-chan int {
	chunkSize := fileSize / chunkCount
	sorts:=[]<-chan int{}
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	for i:=0;i<chunkCount ;i++  {

		file.Seek(int64(i*chunkSize),0)
		source := pipeline.ReaderSource(bufio.NewReader(file), chunkSize)
		source = pipeline.InMeSort(source)
		sorts=append(sorts,source)
	}
	return pipeline.MergeN(sorts...)
}
