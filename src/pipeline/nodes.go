package pipeline

import "sort"

//做一个函数 把数组放到通道里面 只从里面 拿东西
func ArraySource(a ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, v := range a {
			out <- v
		}
		//用完关闭通道  严谨 这里关闭的位置？ 在匿名函数里面和外面的区别？
		close(out)
	}()

	return out
}

/*
在内存中排序
*/
func InMeSort(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		//第一一个切片 但是是不可变的对象  要用append 去收取
		p := []int{}
		for v := range in {
			p = append(p, v)
		}
		sort.Ints(p)
		for _, v := range p {
			out <- v
		}
		close(out)
	}()

	return out
}

//进行合并 第三部
func Merge(in1, in2 <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		v1, ok1 := <-in1
		v2, ok2 := <-in2
		/*
		分析对于两种情况的业务  对于 计算的影响
		*/
		for ok1 || ok2 {
			if !ok2 || (ok1&&v1 <= v2) {
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
