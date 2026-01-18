import com.demo.B;
import com.demo.BMapper;
import com.google.gson.Gson;
import com.sun.prism.impl.Disposer;
import org.apache.ibatis.datasource.pooled.PooledDataSourceFactory;
import org.apache.ibatis.io.Resources;
import org.apache.ibatis.mapping.Environment;
import org.apache.ibatis.session.Configuration;
import org.apache.ibatis.session.SqlSession;
import org.apache.ibatis.session.SqlSessionFactory;
import org.apache.ibatis.session.SqlSessionFactoryBuilder;
import org.apache.ibatis.transaction.TransactionFactory;
import org.apache.ibatis.transaction.jdbc.JdbcTransactionFactory;
import org.apache.ibatis.type.TypeAliasRegistry;

import javax.sql.DataSource;
import java.io.InputStream;

import java.util.List;
import java.util.Properties;

public class Main {


    public static void main(String[] args) {


        try {

            Properties properties = new Properties();
            // Updated driver class name for mysql-connector-j 9.x (replaces com.mysql.jdbc.Driver)
            properties.setProperty("driver", "com.mysql.cj.jdbc.Driver");
            properties.setProperty("url", "jdbc:mysql://127.0.0.1:3306/a");
            properties.setProperty("username", "root");
            properties.setProperty("password", "root");
            PooledDataSourceFactory pooledDataSourceFactory = new PooledDataSourceFactory();
            pooledDataSourceFactory.setProperties(properties);


//            Building SqlSessionFactory without XML
            DataSource dataSource = pooledDataSourceFactory.getDataSource();
            TransactionFactory transactionFactory = new JdbcTransactionFactory();
            Environment environment = new Environment("development", transactionFactory, dataSource);
            Configuration configuration = new Configuration(environment);
            TypeAliasRegistry aliases = configuration.getTypeAliasRegistry();
            aliases.registerAlias("B", B.class);
            configuration.addMapper(BMapper.class);

//成功
            String resource = "mybatisconfig.xml";
            InputStream inputStream = Resources.getResourceAsStream(resource);
            SqlSessionFactory sqlSessionFactory =
                    new SqlSessionFactoryBuilder().build(inputStream);

            //没有成功过
//          SqlSessionFactory  sqlSessionFactory = new SqlSessionFactoryBuilder().build(configuration);

            SqlSession session = sqlSessionFactory.openSession();
            BMapper mapper = session.getMapper(BMapper.class);
            System.out.println("777");
            B entity = mapper.selectB(1);
//            B entity = session.selectOne(
////                    "com.demo.BMapper.selectB", 1);

            Gson gson = new Gson();
            String json = gson.toJson(entity);
            System.out.println("json:" + json);
            System.out.println("-----------------------------------------------");


//            entity.setA("A");
//            entity.setB(666);
//            int rt = session.insert("com.demo.BMapper.insert", entity);
//            if (rt == 0) {
//                System.err.println("Insert into Author failed");
//            }
//            session.commit();
            System.out.println("-----------------------------------------------");
            entity.setB(777);
            session.update("com.demo.BMapper.update", entity);
            session.commit();

            session.delete("com.demo.BMapper.delete", 2);
            session.commit();

            List<B> list = session.selectList("com.demo.BMapper.selectAll");
            System.out.println(gson.toJson(list));
//            session.sele


        } catch (Exception e) {
            System.out.println(e);
        } finally {
            System.out.println("done");
        }

    }


}

