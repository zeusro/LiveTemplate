
## 允许mac安装任意来源的 app

```bash
sudo spctl --master-disable
```

然后在安全性与隐私-通用那里,点击高级,输入密码解锁以允许修改配置,点击任意来源就行

## Alfred3

1. 打开“Alfred 3 KG”，点击“Patch”弹出对话框,找到“Alfred 3.app”（如果弹出安装Xcode命令行工具,安装）
3. 点击“Save”提示“License information saved successfully”完成注册

## Dash + Alfred配置

只需要在安装完了Dash和Alfred后，在Dash的Preferance->Integration选项中点击Alfred

参考链接:
1. [macOS Sierra特性——安全性调整](https://www.feng.com/apple/tutorial/2016-09-27/MacOS-Sierra-features---security-adjustment_658157.shtml)
1. [Mac 神器 Alfred 破解](https://www.jianshu.com/p/72fe06566fce)
1. [Dash + Alfred配置](https://www.jianshu.com/p/77d2bf8df81f)

## Iterm2

    在profiles-Keys那里设置

    右移一个词： opt + 右 send escape sequence f

    左移一个词： opt + 左 send escape sequence b

    cmd+k 清除屏幕
    


参考链接:
1. [开发工具 OS X 下的 iTerm 2 如何让 cursor 跳字移动？](https://ruby-china.org/topics/6114)