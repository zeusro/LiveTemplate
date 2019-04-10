
## 官方包还原

完整安装包见

https://github.com/Microsoft/vscode-go/wiki/Go-tools-that-the-Go-extension-depends-on

`golang.org/x/tools`的依赖包可以用`https://github.com/golang/tools`替换

https://blog.csdn.net/jiankunking/article/details/78544327

## 其他技巧

```go
import(
    "fmt"
    "os"
)

m := `hello
    world`

    
//Go里面有一个关键字iota，这个关键字用来声明enum的时候采用，它默认开始值是0，const中每增加一行加1：
    

a := [...]int{4, 5, 6}

if x := computedValue(); x > 10 {
    fmt.Println("x is greater than 10")
} else {
    fmt.Println("x is less than 10")
}


for _, v := range map{
    fmt.Println("map's val:", v)
}


type testInt func(int) bool // 声明了一个函数类型


cs := make(chan string)

//强制类型断言
v.(type) 

```

## 点操作

    我们有时候会看到如下的方式导入包
```go
import(
    . "fmt"
)
```
    这个点操作的含义就是这个包导入之后在你调用这个包的函数时，你可以省略前缀的包名，也就是前面你调用的fmt.Println("hello world")可以省略的写成Println("hello world")
## 别名操作
    别名操作顾名思义我们可以把包命名成另一个我们用起来容易记忆的名字
```go
import(
    f "fmt"
)
```
    别名操作的话调用包函数时前缀变成了我们的前缀，即f.Println("hello world")
## _操作

这个操作经常是让很多人费解的一个操作符，请看下面这个import
```go
import (
    "database/sql"
    _ "github.com/ziutek/mymysql/godrv"
)
```
_操作其实是引入该包，而不直接使用包里面的函数，而是调用了该包里面的init函数。
