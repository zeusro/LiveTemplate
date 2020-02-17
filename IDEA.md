

## 建包方式

在Intellij IDEA 建立文件夹时总是会一直往后累加，并不是以树的形式直接显示，解决办法是：比如你要在b下面创建一个c,那么你直接添加package是不行的，比如b的上级目录是a，那么你在a上创建一个package里面写上b.c这样就可以在b下面出来c了，而不是c与b平行而形成的b.c的形式.

## idea内置变量

可用于内置方法注释

This is a built-in template. It contains a code fragment that can be included into file templates (Templates tab) with the help of the #parse directive.

The template is editable. Along with static text, code and comments, you can also use predefined variables that will then be expanded like macros into the corresponding values.
Predefined variables will take the following values:

```
${PACKAGE_NAME}	 	name of the package in which the new file is created
${USER}	 	current user system login name
${DATE}	 	current system date
${TIME}	 	current system time
${YEAR}	 	current year
${MONTH}	 	current month
${MONTH_NAME_SHORT}	 	first 3 letters of the current month name. Example: Jan, Feb, etc.
${MONTH_NAME_FULL}	 	full name of the current month. Example: January, February, etc.
${DAY}	 	current day of the month
${HOUR}	 	current hour
${MINUTE}	 	current minute
${PROJECT_NAME}	 	the name of the current project
```

## 常见问题排查

[解决Intellij idea Java JDK多重选择提示问题](https://blog.csdn.net/ruglcc/article/details/72627254)


参考:
1. [IntelliJ IDEA类头注释和方法注释](https://www.oschina.net/question/179541_26961)
2. [IntelliJ IDEA 15 配置 Tomcat8](http://blog.csdn.net/jiankunking/article/details/51921092)
3. [关于The APR based Apache Tomcat Native library警告 - CSDN博客]( http://blog.csdn.net/magerguo/article/details/9467687)
4. [你们都在用IntelliJ IDEA吗？或许你们需要看一下这篇博文](https://www.cnblogs.com/clwydjgs/p/9390488.html)









