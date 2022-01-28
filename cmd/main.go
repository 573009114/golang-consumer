package main

import (
	"golang-consumer/application/controllers"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(4) //CPU核心数
	//创建100个协程
	for i := 0; i < 100; i++ {
		controllers.Consumer()
	}
}
