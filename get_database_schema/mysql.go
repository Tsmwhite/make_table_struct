package get_database_schema

import (
	_ "github.com/go-sql-driver/mysql"
)

type Mysql struct {
	_make_
}

func NewMysql() *Mysql {
	return &Mysql{}
}


func (this *Mysql) GetAllTable(dbName string) {
	sql_str := "SELECT `TABLE_NAME` FROM `information_schema`.`TABLES` WHERE `TABLE_SCHEMA` =?"
	rows, err := this.DbConn.Query(sql_str,dbName)
	if err != nil {
		panic("GetAllTable Query error:"+err.Error())
	}

	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		if err != nil {
			panic("GetAllTable Scan error:"+err.Error())
		}
		this.tableList = append(this.tableList,table)
	}
}

func (this *Mysql) getTableColumns() {
	sql_str := "SELECT `COLUMN_NAME`,`DATA_TYPE`,`TABLE_NAME` FROM `information_schema`.`COLUMNS` WHERE `TABLE_SCHEMA` = ? ORDER BY `TABLE_NAME` ASC,`ORDINAL_POSITION` ASC"
	rows, err := this.DbConn.Query(sql_str,this.DbName)
	if err != nil {
		panic("GetAllTable Query error:"+err.Error())
	}

	var tableColumns  = make(map[string][]column)
	var sortSlice []string
	for rows.Next() {
		var column_name,data_type,table string
		err := rows.Scan(&column_name,&data_type,&table)
		if err != nil {
			panic("GetAllTable Scan error:"+err.Error())
		}
		item := column{
			column_name,
			data_type,
		}
		if len(tableColumns[table]) < 1 {
			sortSlice = append(sortSlice,table)
		}
		tableColumns[table] = append(tableColumns[table],item)
	}

	this.setTableSort(sortSlice)
	this.setTableColumns(tableColumns)
}
