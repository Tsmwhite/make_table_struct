package get_database_schema

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"strings"
)

const(
	_MySql_		= iota
	_SqlServer_
)

type column struct {
	name string
	data_type string
}

type DBOption struct {
	DbName			string
	DbConn			*sql.DB
}

type Maker interface {
	getTableColumns()
	getColumnType(data_type string) string
	setDbOption(*DBOption)
	setTableSort([]string)
	setTableColumns(map[string][]column)
	makeFile()
	makeStructContent()
}

type Make struct {
	DBOption
	connType		int
	tableList		[]string
	tableSort		[]string
	tableColumns	map[string][]column
	fileContent		string
}

var fileDir 	= "./model/"
var fileName 	= "./model/model.go"

func Run(db *DBOption,instance Maker) {
	err := db.DbConn.Ping()
	if err != nil {
		panic("数据库连接无效："+err.Error())
	}

	instance.setDbOption(db)
	instance.getTableColumns()
	instance.makeStructContent()
	instance.makeFile()
}

func (this *Make) setDbOption(db *DBOption) {
	this.DbConn = db.DbConn
	this.DbName = db.DbName
}

func (this Make) makeFile() {
	var file *os.File
	var err  error

	if !checkFileIsExist(fileDir) {
		os.Mkdir(fileDir,os.ModePerm)
	}

	if checkFileIsExist(fileName) { //如果文件存在
		file,err = os.OpenFile(fileName, os.O_WRONLY, os.ModePerm) //打开文件
	} else {
		file,err = os.Create(fileName) //创建文件
	}

	defer file.Close()

	if err != nil {
		panic("文件打开错误："+err.Error())
	}
	_,err = io.WriteString(file,this.fileContent)
	if err != nil {
		panic("文件写入错误："+err.Error())
	}
}

func (this *Make) makeStructContent()  {
	var type_str,structContent string
	var param  []interface{}

	structContent = "package model\n\n"

	for _,table := range this.tableSort {

		param = append(param,HumpFormat(table),)
		structContent += "\ntype %s struct {"

		for _,column := range this.tableColumns[table] {
			type_str = this.getColumnType(column.data_type)
			param = append(param,HumpFormat(column.name),type_str,column.name,)
			structContent += "\n"+
				"  %s %s `json:\"%s\"`"
		}

		structContent += "\n}"
	}
	this.fileContent = fmt.Sprintf(structContent,param...)
}

func (this Make) getColumnType(data_type string) string {
	var res string
	switch this.connType {
		case _MySql_:
			res = mysqlType(data_type)
		case _SqlServer_:
			res = mssqlType(data_type)
	}

	return res
}

func (this *Make) setTableSort(sort []string) {
	this.tableSort = sort
}

func (this *Make) setTableColumns(tableColumns map[string][]column) {
	this.tableColumns = tableColumns
}

func mysqlType(data_type string) string {
	var res string
	switch strings.ToUpper(data_type) {
		case "TINYINT","SMALLINT","MEDIUMINT","INT","BIGINT","INTEGER":
			res = "int64"
		case "FLOAT","DOUBLE","DECIMAL":
			res = "float64"
		default:
			res = "string"
	}
	return  res
}

func mssqlType(data_type string) string {
	var res string
	switch strings.ToUpper(data_type) {
		case "BYTE","INTEGER","LONG","INT","BIGINT":
			res = "int64"
		case "SINGLE","DOUBLE","CURRENCY":
			res = "float64"
		default:
			res = "string"
	}
	return  res
}


func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}


//驼峰格式
func HumpFormat(str string) string {
	var res string
	ary := strings.Split(str,"_")
	for _,v := range ary {
		res += Capitalize(v)
	}
	return res
}

//首字母大写
func Capitalize(str string) string {
	var upperStr string
	vv := []rune(str)
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 {
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				fmt.Println("Not begins with lowercase letter,")
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}
