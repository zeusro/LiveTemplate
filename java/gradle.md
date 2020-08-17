

# Gradle构建Java多项目依赖

## 背景介绍

### 垃圾 maven
使用maven 构建项目有一个很烦人的问题.A项目里面可能用到B依赖, B 依赖于 C,B 依赖 D的1.2版本, C 依赖 D 的1.3版本.

这里面涉及2个原则
1. 路径最近原则（路径最近者优先）
1. 如果传递依赖的路径是一样的，那要看哪个先声明。

    mvn dependency:tree
    
然后生成项目的时候就炸锅了.

## 引入 Gradle

```bash

gradle init

```

## 导入 gradle 项目到 idea

    Do an Import Project or Open... and navigate to build.gradle file.
    
    This should be enough for IntelliJ to figure out the dependencies and set up the project.

    

1. [Maven部署过程中的ClassCastException问题](https://blog.csdn.net/blueheart20/article/details/43448693)
2. [Introduction to the Dependency Mechanism](https://maven.apache.org/guides/introduction/introduction-to-dependency-mechanism.html)
3. [Maven依赖机制](http://fanhongtao.github.io/2013/04/05/maven-dependency-mechanism.html)
4. [Java 构建入门](http://wiki.jikexueyuan.com/project/gradle/java-quickstart.html)
5. [使用Gradle构建Java项目](http://www.importnew.com/15881.html)
6. [gradle配置国内镜像](https://blog.csdn.net/lj402159806/article/details/78422953)
7. [Java 构建入门](http://wiki.jikexueyuan.com/project/gradle/java-quickstart.html)
8. [多项目打包](https://lippiouyang.gitbooks.io/gradle-in-action-cn/content/multi-project/assemble.html)
9. [用 Docker、Gradle 来构建、运行、发布一个 Spring Boot 应用](https://waylau.com/docker-spring-boot-gradle/)
10. [使用Gradle docker插件和registry打包应用| SpringBoot实践](https://www.jianshu.com/p/5b5c034e1da9)
11. [Interface ModuleDependency](https://docs.gradle.org/current/javadoc/org/gradle/api/artifacts/ModuleDependency.html)
12. [Building Java Projects with Gradle](https://spring.io/guides/gs/gradle/)
13. [Creating Multi-project Builds](https://guides.gradle.org/creating-multi-project-builds/?_ga=2.33378940.377288223.1525010261-1262439609.1525010261)
14. [内置变量](https://docs.gradle.org/current/dsl/org.gradle.api.tasks.bundling.Jar.html)
15. []()




