package main

import (
	"bufio"
	"fmt"
	"os"
	"pipeline"
	"time"
)

func main() {
	const filename = "small.txt"
	const n = 64
	start := time.Now().Nanosecond()
	//create a file
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	//想要先生成 一个数据源 然后写入文件
	p := pipeline.RandomSource(n)
	//这里用一个buffer.io来进行缓冲优化速度
	writer := bufio.NewWriter(file)
	pipeline.WriteSink(writer, p)
	/**
	当用到bufio的时候应该注意
	记得一定要进行刷新 不然缓存区里面的 东西写不尽，文件不能完全的保留
	*/
	writer.Flush()
	//从新打开文件读取 数据
	file, err = os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	p = pipeline.ReaderSource(bufio.NewReader(file), -1)
	count := 0
	for v := range p {
		fmt.Println(v)
		count++
		if count >= 10 {
			break
		}
	}
	end := time.Now().Nanosecond()
	fmt.Println("用的时间是：", end-start)
}
func mergeName() {
	p := pipeline.Merge(pipeline.InMeSort(pipeline.ArraySource(2, 3, 1, 7, 6)),
		pipeline.InMeSort(pipeline.ArraySource(4, 3, 8, 23, 11, 44, 9)))
	//for {
	//	if num, ok := <-p; ok {
	//		fmt.Println(num)
	//	} else {
	//		break
	//	}
	//
	//}
	/*用range 来遍历 chan的 数据

	必须也是 有close的关闭通道 要不然 不知道 啥时候关闭
	*/
	for v := range p {
		fmt.Println(v)
	}
}
