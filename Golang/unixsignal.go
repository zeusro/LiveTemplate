package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// var signalChan chan os.Signal

func main() {

	signalChan := make(chan os.Signal)
	fmt.Println("run")
	//使用docker stop 命令去关闭Container时，该命令会发送SIGTERM 命令到Container主进程，让主进程处理该信号，关闭Container，如果在10s内，未关闭容器，Docker Damon会发送SIGKILL 信号将Container关闭。
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	<-signalChan
	signal.Stop(signalChan)
	fmt.Println("I am dead now")
}
