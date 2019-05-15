# 用vs code调试golang项目

## 准备工作
1. 安装vs code(我用的是Windows 1.10.2版),安装对应golang拓展**
1. 假设gopath是D:\GOPATH\src,我们在src下面有2个项目,一个是go-demo,另外一个是a.我们用vs code打开D:\GOPATH\src这个目录(不打开具体项目的那个文件夹是为了方便切换调试的目标)
    

## 单一项目调试
首次打开这个目录的话需要在debug里面添加配置


点击添加配置之后,可以看到src目录下面多了一个.vscode文件夹,.vscode文件夹下面多了launch.json文件
```json
{
    "version": "0.2.0", 
    "configurations": [
        {
            "name": "go-demo", 
            "type": "go", 
            "request": "launch", 
            "mode": "debug", 
            "remotePath": "", 
            "port": 23345, 
            "host": "127.0.0.1", 
            "program": "${workspaceRoot}/go-demo", 
            "env": { }, 
            "args": [ ], 
            "showLog": true
        }, 
        {
            "name": "a", 
            "type": "go", 
            "request": "launch", 
            "mode": "debug", 
            "remotePath": "", 
            "port": 23346, 
            "host": "127.0.0.1", 
            "program": "${workspaceRoot}/a", 
            "env": { }, 
            "args": [ ], 
            "showLog": true
        }
    ]
}
```
按照上面这样配置就行.习惯上,我把name设定的跟项目名一致,项目之间的端口最好不要相同program中的${workspaceRoot}是一个占位符,代表当前打开的D:\GOPATH\src目录.保存之后,会发现原先的"没有配置"多了2个刚才新添加的调试目标,点选一个,点击启动调试,访问http://127.0.0.1:XX就好了


上述的例子适用于每次调试一个项目的情况,如果要同时调试多个项目呢?vs code 也是可以办到的

## 多项目同时调试

在configurations后面,加上
```json
     "compounds": [
        {
            "name": "multi",
            "configurations": ["go-demo", "a"]
        }
    ]
```

可以看到可供调试的项目多了一个multi.


**虽然可以多个项目同时debug,但是调试的时候要分开项目调试**

最终的launch.json
```json
{
    "version": "0.2.0", 
    "configurations": [
        {
            "name": "go-demo", 
            "type": "go", 
            "request": "launch", 
            "mode": "debug", 
            "remotePath": "", 
            "port": 23345, 
            "host": "127.0.0.1", 
            "program": "${workspaceRoot}/go-demo", 
            "env": { }, 
            "args": [ ], 
            "showLog": true
        }, 
        {
            "name": "a", 
            "type": "go", 
            "request": "launch", 
            "mode": "debug", 
            "remotePath": "", 
            "port": 23346, 
            "host": "127.0.0.1", 
            "program": "${workspaceRoot}/a", 
            "env": { }, 
            "args": [ ], 
            "showLog": true
        }
    ], "compounds": [
        {
            "name": "multi",
            "configurations": ["go-demo", "a"]
        }
    ]
}

```

## 断点的技巧

1. 断点只能在debug之前打

本来想写条件断点和命中次数调试的,发现目前golang的插件都不支持,作罢.


## 其他问题

### 启动调试器失败 


Failed to continue: Check the debug console for details.

Failed to continue: "Cannot find Delve debugger. Install from https://github.com/derekparker/delve & ensure it is in your Go tools path, "GOPATH/bin" or "PATH"."

常见于升级golang版本后出现,删了dlv之后重新安装即可

## 参考链接:
1. [Debugging](https://code.visualstudio.com/docs/editor/debugging)
1. []()
