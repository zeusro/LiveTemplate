
    command+\ 分成多栏
    按 F1 就可以调出 Command Palette，这里可以快速输入各种命令，一些基本技巧
    command+D 选择当前单词


## 在Visual Code中使用代码片段复用代码

点击文件-首选项-用户代码片段
![image](resource/QQ截图20161121154103.png)

在下拉中选择go语言

会出现一个json
```json
{
/*
	 // Place your snippets for Go here. Each snippet is defined under a snippet name and has a prefix, body and 
	 // description. The prefix is what is used to trigger the snippet and the body will be expanded and inserted. Possible variables are:
	 // $1, $2 for tab stops, ${id} and ${id:label} and ${1:label} for variables. Variables with the same id are connected.
	 // Example:
	 "Print to console": {
		"prefix": "log",
		"body": [
			"console.log('$1');",
			"$2"
		],
		"description": "Log output to console"
	}
*/
}
```

注释里面是默认的示例,我们可以根据需要定义一个适合自己的代码片段


```json

"生成API方法": {
		"prefix": "api",
		"body": [
			"type ${apiname:输入API名称}API struct {",
			"	apiserver.Controller",
			"}",
			"",
			"func (a * ${apiname}API) Prepare() {",
			"",
			"}",
			"",
			"func init() {",
			"	apiserver.AddController(\"${action}\", &${apiname}API{})",
			"",
			"}",
			"",
			"func (a *${apiname}API) FuncMapping() {",
			"",
			"",
			"}",
			"",
			"",
			""
		],
		"description": "生成API方法"
	}
```

保存之后,在*.go文件中键入api-tab就会出现并导出这个模板${}是变量名,在我定义的模板里面,一共有2个变量.输入第一个变量之后,再按tab就可以切换到${action}那个位置

etc


参考链接:

1. [Adding Snippets to Visual Studio Code](https://code.visualstudio.com/Docs/customization/userdefinedsnippets)
1. []()
1. []()
1. []()
1. []()