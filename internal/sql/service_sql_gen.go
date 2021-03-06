package sql

import (
	"github.com/poemp/goway/inter"
	"github.com/poemp/goway/internal/entity"
	"reflect"
	"strings"
	"unicode"
)

// 判断是否大写
func upper(s string) bool {
	return unicode.IsUpper([]rune(s)[0])
}

// 格式话数据列
// 当大写的字母时候， 转化成小写字母， 并且在前面加上下划线
func formatField(field string) string {
	var tmp string
	for i, k := range field {
		up := false
		// 小写
		if upper(string(string(k))) {
			// 转化成大写
			k += 32
			up = true
		}
		// 非第一位， 并且大写， 转化成小写过后， 在前面添加下划线
		if up && i != 0 {
			tmp += "_"
		}
		tmp += string(k)

	}
	return tmp
}

// 使用反射， 拿到所有的列
// 返回数据列
func GetAllField() []entity.FieldType {
	var history entity.SchemaHistory
	var fields []entity.FieldType
	t := reflect.TypeOf(history)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fields = append(fields, entity.FieldType{
			TableColumn: formatField(f.Name),
			FieldName:   f.Name,
			Kind:        f.Type.Kind(),
		})
	}
	return fields
}

// 获取初始化数据表
func GetTableInitSql() string {
	fieldTypes := GetAllField()
	table := inter.DefaultTableDataSource()
	sql := "create table  " + table.SearchPath + "." + table.Table + "(\n"
	var fields []string
	for _, f := range fieldTypes {
		switch f.Kind {
		//如果是string
		case reflect.String:
			{
				fields = append(fields, "	"+f.TableColumn+"  varchar(200)		null ")
			}
		case reflect.Int64:
			{
				fields = append(fields, "	"+f.TableColumn+"  int8		null ")
			}
		case reflect.Bool:
			{
				fields = append(fields, "	"+f.TableColumn+"  bool 	null")
			}
		default:
			fields = append(fields, "	"+f.TableColumn+"  TIMESTAMP")
		}
	}
	sql += strings.Join(fields, ",\n")
	sql += "\n);"
	return sql
}

// 获得查询sql
func GetSelectSQL() string {
	fieldTypes := GetAllField()
	table := inter.DefaultTableDataSource()
	sql := "   "
	var fields []string
	for _, f := range fieldTypes {
		fields = append(fields, f.TableColumn)
	}
	sql += strings.Join(fields, ", ")
	sql += " from " + table.SearchPath + "." + table.Table + " order by installed_on ;"
	return sql
}

// 获取表是否存在
func GetExistSQL() string {
	table := inter.DefaultTableDataSource()
	sql := "select 1 as col   from pg_tables where schemaname = '" + table.SearchPath + "' and tablename = '" + table.Table + "'"
	return sql
}

// 获取插入数据的sql
func GetInsertIntoSql(values []string) string {

	fieldTypes := GetAllField()
	table := inter.DefaultTableDataSource()
	sql := "insert into  " + table.SearchPath + "." + table.Table
	var fields []string
	for _, f := range fieldTypes {
		fields = append(fields, f.TableColumn)
	}
	sql += " (" + strings.Join(fields, ", ") + " ) "
	sql += " values ('" + strings.Join(values, "','") + "');"
	return sql
}
