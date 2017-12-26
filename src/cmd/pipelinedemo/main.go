package main

import (
	"fmt"
	"pipeline"
)

func main() {
	p :=pipeline.Merge(pipeline.InMeSort( pipeline.ArraySource(2, 3, 1, 7, 6)),
		pipeline.InMeSort( pipeline.ArraySource(4,3,8,23,11,44,9)))
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
