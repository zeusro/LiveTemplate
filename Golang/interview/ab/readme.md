# ab - 两协程交替顺序打印

两个 Go 协程交替顺序打印 1..n（一个打奇数，一个打偶数），三种实现方式。

## 运行

在仓库根目录（含 `go.mod`）下：

```bash
# 默认：打印 1..10，方式 two_channels
go run ./Golang/interview/ab/cmd

# 指定 n 和方式
go run ./Golang/interview/ab/cmd 6 two
go run ./Golang/interview/ab/cmd 6 single
go run ./Golang/interview/ab/cmd 6 cond
```

## 实现说明

| 文件 | 方式 | 思路 |
|------|------|------|
| `two_channels.go` | 两个无缓冲 channel | 各持有一个 channel，打印完后把“令牌”发给对方，自己再等对方的令牌 |
| `single_channel.go` | 单一 channel | channel 传 `(数值, 轮到谁)`；每个协程只处理自己的消息，否则放回 channel，保证严格交替 |
| `cond.go` | Mutex + Cond | 共享 `turn`（0=奇数方，1=偶数方），不是自己回合就 `Wait`，打完改 `turn` 并 `Signal` |

## 示例输出

```
交替打印 1..10，方式: two

1
2
3
4
5
6
7
8
9
10
```
