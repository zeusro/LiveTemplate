# 类型转换

## array tostring

  s := string(byteArray[:])


## int64.int32... tostring
 strconv.FormatInt(int64(int64), 10)
 int tostring
 strconv.Itoa()

## float64 tostring
 strconv.FormatFloat(i.(float64), 'f', -1, 64)

## int to int64
 int64(xx)

## string to byte

 strconv.ParseInt(string, 10, 8)

## string to int
 int,err:=strconv.Atoi(string)

## string to int32

 strconv.ParseInt(v32, 10, 32)

## string to int64
 int64Value, err := strconv.ParseInt(string, 10, 64)

## 全角半角相互转换

```go
import "golang.org/x/text/width"

s := `。，（）-1！@234567890abc１２３４５６７８９ａｂｃ`
// 全角转半角
fmt.Println(width.Narrow.String(s))
// 半角转全角
fmt.Println(width.Widen.String(s))
// 参考 https://www.golangnote.com/topic/255.html
```
