package com.demo;

import org.apache.ibatis.annotations.Param;
import org.apache.ibatis.annotations.Select;

public interface BMapper {


//    @Select("SELECT * FROM B WHERE id = #{id}")
//    B selectBWithAnnotations(Integer id);

    B selectB(@Param("id") Integer id);

//    B selectBWithAuthor(Integer id);
}
