# make_table_struct
go 连接数据库自动生成struct

```
go get github.com/theSmallwhiteMe/make_table_struct

新版本
go get github.com/Tsmwhite/structMaker
```


```golang
package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/theSmallwhiteMe/make_table_struct/get_database_schema"
)

func main() {
	dbConfigString := "root:lxy196914@tcp(127.0.0.1:3306)/white_blog?charset=utf8"
	db,err := sql.Open("mysql",dbConfigString)
	if err == nil {
		dbOption := &get_database_schema.DBOption{
			"white_blog",
			db,
		}
		get_database_schema.Run(dbOption,get_database_schema.NewMysql())
	}
}
```

新手上路，请大家多多指教
