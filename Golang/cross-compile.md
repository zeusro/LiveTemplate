
## golang交叉编译技巧

### windows
- 切换到window64编译环境
```
set GOARCH=amd64
set GOOS=windows
go tool dist install -v pkg/runtime
go install -v -a std
go build
```

- 切换到linux64编译环境
```
set GOARCH=amd64
set GOOS=linux
go tool dist install -v pkg/runtime
go install -v -a std
go build  

```

**切换是永久切换.比如切换到linux64位环境之后,之后的go build就都生成的linux64位二进制文件**

### unix / linux

https://github.com/golang/go/wiki/WindowsCrossCompiling




    utf8.RuneCountInString(str)
    go list -json
    go.exe test -timeout 30s -v -run ^Test_corpNotify_infom_list$

    

参考链接:

1. [How to cross compile from Windows to Linux?](http://stackoverflow.com/questions/20829155/how-to-cross-compile-from-windows-to-linux)
1. [go install](http://wiki.jikexueyuan.com/project/go-command-tutorial/0.2.html)
1. [WindowsBuild](https://github.com/golang/go/wiki/WindowsBuild)
1. [WindowsCrossCompiling](https://github.com/golang/go/wiki/WindowsCrossCompiling)