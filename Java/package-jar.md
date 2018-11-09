# 把java项目打包成一个可执行jar

### maven项目

```xml
<build>
  <plugins>
    <plugin>
      <artifactId>maven-assembly-plugin</artifactId>
      <configuration>
        <archive>
          <manifest>
            <mainClass>fully.qualified.MainClass</mainClass>
          </manifest>
        </archive>
        <descriptorRefs>
          <descriptorRef>jar-with-dependencies</descriptorRef>
        </descriptorRefs>
      </configuration>
    </plugin>
  </plugins>
</build>
```
```
mvn clean compile assembly:single
```



参考链接:
1. [How to Create an Executable JAR with Maven](http://www.baeldung.com/executable-jar-with-maven)
1. [How can I create an executable JAR with dependencies using Maven?](https://stackoverflow.com/questions/574594/how-can-i-create-an-executable-jar-with-dependencies-using-maven?answertab=active#tab-top)
1. []()
1. []()
1. []()
1. []()
1. []()
1. []()
