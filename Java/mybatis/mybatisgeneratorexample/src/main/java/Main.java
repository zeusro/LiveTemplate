import com.demo.dao.ActorMapper;
import com.demo.entity.Actor;
import com.demo.entity.ActorExample;
import com.google.gson.Gson;
import jdk.nashorn.internal.objects.Global;
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
import org.apache.log4j.Logger;

import javax.sql.DataSource;
import java.io.InputStream;
import java.nio.file.ClosedWatchServiceException;
import java.text.SimpleDateFormat;
import java.util.*;

import static java.lang.System.out;
import static java.util.Calendar.FEBRUARY;
import static java.util.Calendar.JANUARY;

public class Main {


    final static Logger logger = Logger.getLogger(Main.class);

    public static void main(String[] args) {
        try {


            String resource = "mybatisconfig.xml";
            InputStream inputStream = Resources.getResourceAsStream(resource);
            SqlSessionFactory sqlSessionFactory = new SqlSessionFactoryBuilder().build(inputStream);
            Configuration configuration = sqlSessionFactory.getConfiguration();
            SqlSession session = sqlSessionFactory.openSession();
            ActorMapper mapper = session.getMapper(ActorMapper.class);
            Actor entity;
            List<Actor> entityList;

            entity = mapper.selectByPrimaryKey(Short.valueOf("1"));
            Gson gson = new Gson();
            out.println(gson.toJson(entity));

//            entityList = mapper.selectByExample(ConditionSelect());
            //            out.println(gson.toJson(entityList));

            Long ddd = mapper.countByExample(ConditionSelect());
            out.println("ddd:" + ddd);

//            String parameter = "parameter";
//            logger.debug("This is debug : " + parameter);
//            logger.info("This is info : " + parameter);
//            logger.warn("This is warn : " + parameter);
//
//            logger.error("This is error : " + parameter);
//            logger.fatal("This is fatal : " + parameter);


        } catch (Exception e) {
            out.println(e);
        } finally {
            out.println("done");
        }
        out.println("-----------------------------------------------");


    }

    public static ActorExample ConditionSelect() {
        ActorExample example = new ActorExample();

        example.or().andActorIdGreaterThan((short) 1).andActorIdLessThan((short) 3);
//        SimpleDateFormat isoFormat = new SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss");
//        isoFormat.setTimeZone(TimeZone.getTimeZone("GMT+08:00"));

//        Date date1 = new GregorianCalendar(2006, FEBRUARY, 9).getTime();
//        Date date2 = new GregorianCalendar(2006, FEBRUARY, 19).getTime();
//        example.or().andLastUpdateBetween(date1, date2);
        example.or().andFirstNameLike("%E%");
        return example;
    }

}

