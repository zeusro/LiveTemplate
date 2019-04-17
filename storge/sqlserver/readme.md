
1. [实用的查询](#实用的查询)
1. [解决sqlserver单用户模式无法改成多用户模式](#解决sqlserver单用户模式无法改成多用户模式)
1. [每日增量备份](#每日增量备份)




## 实用的查询

```
select * from sysobjects where xtype='u' ORDER BY name
select *  from sysobjects where xtype='V' ORDER BY name
select * from sysobjects where xtype='P' ORDER BY name

sql server 数表:
select count(1) from sysobjects where xtype='U'
数视图:
select count(1) from sysobjects where xtype='V'
数存储过程
select count(1) from sysobjects where xtype='P'
```

参考链接:

1. [SQL Server中，查询数据库中有多少个表，以及数据库其余类型数据统计查询](http://blog.csdn.net/u012138032/article/details/78455990)




## 解决sqlserver单用户模式无法改成多用户模式


先要让数据库脱机

```
ALTER DATABASE <DATABASE NAME> SET AUTO_UPDATE_STATISTICS_ASYNC OFF
```


```sql
declare @kid varchar(8000)
set @kid=''
select @kid=@kid+' kill '+cast(spid as varchar(8))
from master..sysprocesses
where dbid=db_id('数据库名')
Exec(@kid)
```


```sql
-- 设置数据库立马脱机
ALTER DATABASE dbname SET OFFLINE WITH ROLLBACK IMMEDIATE 
-- 立马上线
alter database dbname set  online  
```

参考链接:
1. [alter table set](https://docs.microsoft.com/zh-cn/sql/t-sql/statements/alter-database-transact-sql-set-options)
2. [将数据库设置为单用户模式](https://docs.microsoft.com/zh-cn/sql/relational-databases/databases/set-a-database-to-single-user-mode)
3. [sqlserver单用户模式无法改成多用户模式](http://www.wangjingfeng.com/491.html)
4. [sqlserver独占数据库](http://www.bkjia.com/Sql_Server/490634.html)
1. [Creating, detaching, re-attaching, and fixing a SUSPECT database](https://www.sqlskills.com/blogs/paul/creating-detaching-re-attaching-and-fixing-a-suspect-database/)


## 每日增量备份

用管理员 Windows 或者 sa 登录服务器,在作业那里傻瓜式添加新的作业即可,重点的 sql 是下面这句

```sql
DECLARE @name varchar(50)
DECLARE @datetime varchar(14)
DECLARE @path varchar(255)
DECLARE @bakfile varchar(255)

set @name='<dbname>'
set @datetime=CONVERT(char(8),getdate(),112)
set @path='d:\backup'
set @bakfile=@path+'\'+@name+'\'+@datetime+'.BAK'
backup database @name to DISK = @bakfile
with COMPRESSION,CHECKSUM,NOFORMAT,NOINIT,SKIP, REWIND, NOUNLOAD,STATS = 10;
GO
```

* 用递归的公用表达式获取父节点

```sql
  with Cte_A([fileid],[parentFileId])
  AS(
  SELECT [fileid],[parentFileId] FROM [a].[dbo].[FileInfo]  where [fileid]=6  union all 
  SELECT a.[fileid],a.[parentFileId] FROM [a].[dbo].[FileInfo] a  inner join Cte_A c on a.fileid=c.[parentFileId]
  )
  select * from Cte_A 
OPTION (MAXRECURSION 2)
```

* 用公用表达式获取所有子节点

```sql
  with Cte_A([fileid],[parentFileId])
  AS(
  SELECT [fileid],[parentFileId] FROM [a].[dbo].[FileInfo]  where [fileid]=1  union all 
  SELECT a.[fileid],a.[parentFileId] FROM [a].[dbo].[FileInfo] a  inner join Cte_A c on a.[parentFileId]=c.fileid
  )
  select * from Cte_A 
OPTION (MAXRECURSION 2)
```
