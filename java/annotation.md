# spring相关注解

todo:

@RequestHeader(value="User-Agent")


## 配置类

### @Profile("production")


### @PropertySource

    @PropertySource(value = "classpath:application-${env}.properties", encoding = "UTF-8")

显式导入配置

### @Value("${property}")

用于读取配置,但是遇到`.properties`文件时会出问题,因为读取文件用的`ISO 8859-1`编码.解决办法可以是
1. 使用`yml`配置(但是`@PropertySource`将不可用)
2. idea设置转换为ASCⅡ
3. 用注解`@PropertySource(value = "classpath:application-${env}.properties", encoding = "UTF-8")`的方式导入

参考:
1. [UTF-8 encoding of application.properties attributes in Spring-Boot](https://stackoverflow.com/questions/37436927/utf-8-encoding-of-application-properties-attributes-in-spring-boot)
2. 

### @ConfigurationProperties

    @ConfigurationProperties(prefix = ZWD_PREFIX) 

将特定配置转化为实体
    
参考:
1. [spring boot 使用@ConfigurationProperties](https://blog.csdn.net/yingxiake/article/details/51263071)

### @ControllerAdvice

将 Controller 层的异常和数据校验的异常进行统一处理，减少模板代码，减少编码量，提升扩展性和可维护性。

参考:
1. [@ControllerAdvice + @ExceptionHandler 全局处理 Controller 层异常](https://blog.csdn.net/kinginblue/article/details/70186586)
2. 


## 数据类


### @DataJpaTest/@JdbcTest/@DataMongoTest

用来测试数据库的,默认使用一个嵌入式内存数据库.

```java
@RunWith(SpringRunner.class)
@DataJpaTest
public class ExampleRepositoryTests {

    @Autowired
    private TestEntityManager entityManager;

    @Autowired
    private UserRepository repository;

    @Test
    public void testExample() throws Exception {
        this.entityManager.persist(new User("sboot", "1234"));
        User user = this.repository.findByUsername("sboot");
        assertThat(user.getUsername()).isEqualTo("sboot");
        assertThat(user.getVin()).isEqualTo("1234");
    }
}

@RunWith(SpringRunner.class)
@JdbcTest
@Transactional(propagation = Propagation.NOT_SUPPORTED)
public class ExampleNonTransactionalTests {

}
```

参考:
1. [透彻的掌握 Spring 中@transactional 的使用](https://www.ibm.com/developerworks/cn/java/j-master-spring-transactional-use/index.html)

### @Cacheable

```java
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Component;

@Component
public class MathService {

	@Cacheable("piDecimals")
	public int computePiDecimal(int i) {
		// ...
	}

}
```
@Cacheable可以标记在一个方法上，也可以标记在一个类上。当标记在一个方法上时表示该方法是支持缓存的，当标记在一个类上时则表示该类所有的方法都是支持缓存的。对于一个支持缓存的方法，Spring会在其被调用后将其返回值缓存起来，以保证下次利用同样的参数来执行该方法时可以直接从缓存中获取结果，而不需要再次执行该方法。Spring在缓存方法的返回值时是以键值对进行缓存的，值就是方法的返回结果，至于键的话，Spring又支持两种策略，默认策略和自定义策略，这个稍后会进行说明。需要注意的是当一个支持缓存的方法在对象内部被调用时是不会触发缓存功能的。@Cacheable可以指定三个属性，value、key和condition。

参考:
1. [Spring缓存注解@Cacheable、@CacheEvict、@CachePut使用](https://blog.csdn.net/wjacketcn/article/details/50945887)