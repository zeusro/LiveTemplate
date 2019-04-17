
## PowerShell


```bash
获取文件md5
Get-FileHash "D:\software\cn_windows_10_multi-edition_version_1709_updated_nov_2017_x64_dvd_100290206.iso" -Algorithm MD5
Get-FileHash "D:\software\cn_windows_10_multi-edition_version_1709_updated_sept_2017_x64_dvd_100090804.iso" -Algorithm SHA1
```


## 项识别为 cmdlet、函数、脚本文件或可运行程序的名称。请检查名称的拼写，如果包括路径，请确保路径正确，然后再试一次

用户如果是第一次使用powershell 执行脚本 的话。其中的原因是：
windows默认不允许任何脚本运行，你可以使用"Set-ExecutionPolicy"cmdlet来改变的你PowerShell环境。

你可以使用如下命令让PowerShell运行在无限制的环境之下：

    Set-ExecutionPolicy Unrestricted


## 获取文件SHA1值
```ps
Get-FileHash "D:\software\cn_windows_10_multi-edition_version_1709_updated_sept_2017_x64_dvd_100090804.iso" -Algorithm SHA1
```

## 获取文件MD5值

```ps
Get-FileHash "D:\software\cn_windows_10_multi-edition_version_1709_updated_nov_2017_x64_dvd_100290206.iso" -Algorithm MD5
Get-FileHash "D:\software\cn_windows_10_multi-edition_version_1709_updated_sept_2017_x64_dvd_100090804.iso" -Algorithm SHA1
```