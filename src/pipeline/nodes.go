package pipeline

import (
	"encoding/binary"
	"io"
	"math/rand"
	"sort"
	"time"
	"fmt"
)
var startime time.Time
func Init(){
	startime=time.Now()
}
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
	out := make(chan int,1024)
	go func() {
		//第一一个切片 但是是不可变的对象  要用append 去收取
		p := []int{}
		for v := range in {
			p = append(p, v)
		}
		fmt.Println("read done:::",time.Now().Sub(startime))
		//运用 标准库中的sort进行排序
		sort.Ints(p)
		fmt.Println("sort done:::",time.Now().Sub(startime))
		//排序之后再放进通道里面
		for _, v := range p {
			out <- v
		}
		close(out)
	}()

	return out
}

//进行合并 第三部 核心的算法
func Merge(in1, in2 <-chan int) <-chan int {
	out := make(chan int,1024)
	go func() {
		v1, ok1 := <-in1
		v2, ok2 := <-in2
		/*
			分析对于两种情况的业务  对于 计算的影响
		*/
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
		fmt.Println("Merge done:::",time.Now().Sub(startime))
	}()
	return out
}

//读文件  改造方法进行分块读取  块的尺寸是-1的时候全部读取
func ReaderSource(reader io.Reader, chuakSize int) <-chan int {
	out := make(chan int, 8)
	go func() {
		buffer := make([]byte, 8)
		bytesread := 0
		for {
			n, err := reader.Read(buffer)
			bytesread += n
			if n > 0 {
				//当n>0是说明 还有 数据要读取  把buffer中的 数据放入通道
				//binary.BigEndian是 二进制工具类的大端实现 int等于cpu的位数与unint等是 兄弟类型？
				v := int(binary.BigEndian.Uint64(buffer))
				out <- v
			}
			if err != nil || (chuakSize != -1 && bytesread >=chuakSize) {
				break
			}
		}
		close(out)
	}()
	return out
}

//写数据进入 通道
func WriteSink(writer io.Writer, in <-chan int) {
	for v := range in {
		buffer := make([]byte, 8)
		binary.BigEndian.PutUint64(buffer, uint64(v))
		writer.Write(buffer)
	}
}

//随机生成数字的 资源
func RandomSource(count int) <-chan int {
	out := make(chan int,1024)
	go func() {
		for i := 0; i < count; i++ {
			out <- rand.Int()
		}
		close(out)
	}()
	return out
}

//节点之间的归并 两两归并
func MergeN(inputs ...<-chan int) <-chan int {
	if len(inputs) == 1 {
		//说明 只有一个元素
		return inputs[0]
	}
	m := len(inputs) / 2
	//merge inputs[0..m) and inputs(m..end]
	return Merge(MergeN(inputs[0:m]...),
		MergeN(inputs[m:]...))
}
