package main

import (
	"fmt"
	"net/http"
)

func main() {
	//实现 对网页的链接
	http.HandleFunc("/", func(writer http.ResponseWriter, rq *http.Request) {
		fmt.Fprintf(writer, "<h1>hello server!!!%s</h1>", rq.FormValue("name"))
	})
	fmt.Println("test debug!")
	//默认是 localhost 的开始
	http.ListenAndServe(":8000", nil)
}
