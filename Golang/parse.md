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
