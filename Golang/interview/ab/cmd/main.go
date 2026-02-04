package main

import (
	"fmt"
	"os"
	"strconv"

	"interview/ab"
)

func main() {
	n := 10
	if len(os.Args) > 1 {
		if v, err := strconv.Atoi(os.Args[1]); err == nil && v > 0 {
			n = v
		}
	}

	mode := "two"
	if len(os.Args) > 2 {
		mode = os.Args[2]
	}

	fmt.Printf("交替打印 1..%d，方式: %s\n\n", n, mode)

	switch mode {
	case "two", "two_channels":
		ab.RunTwoChannels(n)
	case "single", "single_channel":
		ab.RunSingleChannel(n)
	case "cond":
		ab.RunCond(n)
	default:
		fmt.Println("用法: go run . [n] [two|single|cond]")
		fmt.Println("  n: 打印到 n，默认 10")
		fmt.Println("  two_channels: 两个无缓冲 channel（默认）")
		fmt.Println("  single_channel: 单一 channel 传 (值, 轮到谁)")
		fmt.Println("  cond: Mutex + Cond")
	}
}
