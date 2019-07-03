# fmt

* [fmt.Sprintf](#fmt.Sprintf)
* [控制台打印各种ANSI转义序列](#color)


## fmt.Sprintf

可以通过%[n]X达到参数复用的效果,如果待格式化的字符串含有%,可以使用%%

```go
	sql := fmt.Sprintf("%[1]s %[2]s %[3]d %[1]s %%", "1", "2", 3)
	fmt.Println(sql)
// 1 2 3 1 %
```

## color

```go
	flag := 0
	str := "Zeusro "
	for r := 255; r >= 0; r -= 2 {
		for g, b := 0, 255; g < 255 && b >= 1; g += 1 {
			if flag >= len(str) {
				flag = 0
			}
			a := str[flag]
			flag++
			fmt.Printf("\x1b[48;2;%s;%s;%sm\x1b[38m%s\x1b[0m", strconv.Itoa(r), strconv.Itoa(g), strconv.Itoa(b), string(a))
			b -= 1
		}
		fmt.Printf("\x1bE")
	}
```	

[ANSI 终端输出瞎搞指北](https://learnku.com/articles/26231?order_by=vote_count&)

参考链接:

1. []()
