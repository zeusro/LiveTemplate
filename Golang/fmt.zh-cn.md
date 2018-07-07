# fmt

* [fmt.Sprintf](#fmt.Sprintf)
* [](#)


## fmt.Sprintf

可以通过%[n]X达到参数复用的效果,如果待格式化的字符串含有%,可以使用%%

```go
	sql := fmt.Sprintf("%[1]s %[2]s %[3]d %[1]s %%", "1", "2", 3)
	fmt.Println(sql)
// 1 2 3 1 %
```


参考链接:

1. []()
