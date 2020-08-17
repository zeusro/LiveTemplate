
## 单列索引

  CREATE INDEX `idx_`  ON `daifa`.`check_result` (created_at) COMMENT '' ALGORITHM DEFAULT LOCK DEFAULT;

## 组合索引

  CREATE INDEX `mult_`  ON `shoppe_img`.`gms_album4desc` (imgs_width, imgs_height) COMMENT '' ALGORITHM DEFAULT LOCK DEFAULT

  CREATE INDEX `idx_hot_level`  ON  `17zwdv3_product`.z_goods_info(hot_level) COMMENT '' ALGORITHM DEFAULT LOCK DEFAULT;

## 根据表名查数据库

  select TABLE_SCHEMA from information_schema.tables where table_name ='table_name';

##  mysql Client Options

  mysql --user=xxx  --password=xxxxx  -h xxxxx.mysql.rds.aliyuncs.com

https://dev.mysql.com/doc/refman/8.0/en/mysql-command-options.html
