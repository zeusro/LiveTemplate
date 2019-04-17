# 7-zip命令行


## 全覆盖解压
```
"7z.exe" x "X.7z"   -o"导出目录"   -y
```


## 查看版本
```
7z -version
```

## 压缩
```
7z a "fuckyou.7z" "目标目录"  -r-
7z a -tzip archive.gz @listfile.txt
7z a "fuckyo2.7z" "目标目录"  -r0
7z a "fuckyo2.7z" "目标目录" -r
"7z.exe"  a  "E:\Files.7z" "文件1路径" "文件2路径" -r
```



## 参考链接
1. [https://sevenzip.osdn.jp/chm/cmdline/syntax.htm](https://sevenzip.osdn.jp/chm/cmdline/syntax.htm)
1. [https://sevenzip.osdn.jp/chm/cmdline/switches/index.htm](https://sevenzip.osdn.jp/chm/cmdline/switches/index.htm)
