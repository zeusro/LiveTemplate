# 学习mybatis


* [最简单的实例:完全手写](#最简单的实例:完全手写)
* []()
* []()
* []()
* []()





## 最简单的实例:完全手写




参考链接:
1. [Mybatis中接口和对应的mapper文件位置配置深入剖析](http://blog.csdn.net/lmy86263/article/details/53428417)
1. [MyBatis传入多个参数的问题 - mingyue1818 - 博客园 ](http://www.cnblogs.com/mingyue1818/p/3714162.html)
1. [Java Code Examples for org.apache.ibatis.session.SqlSession.getMapper()](https://www.programcreek.com/java-api-examples/index.php?class=org.apache.ibatis.session.SqlSession&method=getMapper)
1. [How do I build SqlSessionFactory without XML? | Kode Java ](https://kodejava.org/how-do-i-build-sqlsessionfactory-without-xml/)
1. [Getting Started with MyBatis 3: CRUD Operations Example with XML Mapper ](https://www.concretepage.com/mybatis-3/getting-started-with-mybatis-3-crud-operations-example-with-xml-mapper)
1. [MyBatis tutorial - Introductory MyBatis tutorial]( http://zetcode.com/db/mybatis/)
1. []()
1. []()
1. []()
1. []()



## 使用 MyBatis Generator快速生成XML和响应mapper以及example

参考[mybatis-generator 代码自动生成工具（maven方式）](http://www.cnblogs.com/JsonShare/p/5521901.html)即可,这个没什么难度

我的示例项目:
[mybatisgeneratorexample](https://github.com/zeusro/LiveTemplate/tree/master/Java/mybatis/mybatisgeneratorexample)


* pom.xml-dependencies

```xml
 <!-- https://mvnrepository.com/artifact/mysql/mysql-connector-java -->
        <dependency>
            <groupId>mysql</groupId>
            <artifactId>mysql-connector-java</artifactId>
            <version>6.0.6</version>
        </dependency>

        <!-- https://mvnrepository.com/artifact/org.mybatis/mybatis -->
        <dependency>
            <groupId>org.mybatis</groupId>
            <artifactId>mybatis</artifactId>
            <version>3.4.5</version>
        </dependency>
```

可能会遇到


1. [Exception getting JDBC Driver](https://blog.csdn.net/u012995964/article/details/53887534)

* pom.xml-build-plugins
```xml
 <plugin>
                <groupId>org.mybatis.generator</groupId>
                <artifactId>mybatis-generator-maven-plugin</artifactId>
                <version>1.3.5</version>

                <executions>
                    <execution>
                        <id>Generate MyBatis Artifacts</id>
                        <goals>
                            <goal>generate</goal>

                        </goals>

                    </execution>
                </executions>
            </plugin>
```

参考链接:
1. [菜鸟在MyBatis路上前行-MyBatis Generator快速生成实体类代码 | Souvc]( http://www.souvc.com/?p=2811)
1. [Example Class Usage Notes](http://www.mybatis.org/generator/generatedobjects/exampleClassUsage.html)
1. [[mybatis]Example的用法](http://blog.csdn.net/zhemeban/article/details/71901759)
1. [mybatis-generator 代码自动生成工具（maven方式）](http://www.cnblogs.com/JsonShare/p/5521901.html)
1. [使用Mybatis-Generator自动生成Dao、Model、Mapping相关文件](http://www.cnblogs.com/lichenwei/p/4145696.html)
1. [Intellij IDEA 14中使用MyBatis-generator 自动生成MyBatis代码](http://blog.csdn.net/luanlouis/article/details/43192131)
1. []()
1. []()
1. []()


## 通过springboot使用mybatis

参考链接:
1. [springboot(六)：如何优雅的使用mybatis - ityouknow's Blog ](http://www.ityouknow.com/springboot/2016/11/06/springboot(%E5%85%AD)-%E5%A6%82%E4%BD%95%E4%BC%98%E9%9B%85%E7%9A%84%E4%BD%BF%E7%94%A8mybatis.html)
1. []()
1. []()
1. []()
1. []()
1. []()
1. []()
1. []()
1. []()

## 其他mybatis技巧

### 控制台打印日志

如果你用SpringBoot，那么在application.properties文件加入以下语句就好了：logging.level.com.xxx.xxx.*.mapper=debug（com开始，就是你mapper接口所在的包路劲）号外：SpringBoot真是好东西啊。

参考链接:
1. [mybatis学习之路----打印sql语句](http://blog.csdn.net/xu1916659422/article/details/78093108)
1. []()
1. []()
1. []()
1. []()
1. []()
1. []()
1. []()


## 疑难解答/排错


?characterEncoding=utf8


参考链接:
1. [Table configuration with catalog null, schema null错误的一个原因 - CSDN博客 ](http://blog.csdn.net/yewei11/article/details/41929179)
1. [mybatis错误——java.io.IOException: Could not find resource com/xxx/xxxMapper.xml - 阿飞(dufyun)的博客 - CSDN博客]( http://blog.csdn.net/u010648555/article/details/70880425)
1. []()
1. []()
1. []()
1. []()
1. []()
1. []()





其他参考链接:
1. [SQL语句构建器类](http://www.mybatis.org/mybatis-3/zh/statement-builders.html)
1. []()
1. []()
1. []()
1. []()
1. []()

