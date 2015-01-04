package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	m "github.com/pravasan/pravasan/migration"

	_ "github.com/lib/pq"
)

var Db *sql.DB
var working_version string
var local_updown string

func init() {
	// fmt.Println("postgres init() it runs before other functions")
}

func Init(c m.Config) {
	fmt.Println("Inside the Postgres")
	Db, _ = sql.Open("mysql", "root:root@/onetest")
}

func GetLastMigrationNo() string {
	var max_version string = ""
	query := "SELECT max(`version`) FROM `schema_migrations`"
	q, err := Db.Query(query)
	defer q.Close()
	if err != nil {
		log.Println("schema_migrations table doesn't exists")
		log.Fatal(err)
	} else {
		q.Next()
		q.Scan(&max_version)
	}
	return max_version
}

func CreateMigrationTable() {
}

func ProcessNow(lm m.Migration, mig m.UpDown, updown string) {
	if updown == "up" && lm.Id <= GetLastMigrationNo() {
		return
	}
	local_updown = updown

	working_version = lm.Id
	nid, _ := strconv.Atoi(lm.Id)
	if nid != 0 {
		fmt.Println("ID : ", lm.Id)

		for _, v := range mig.Create_Table {
			var values_array []string
			for _, vv := range v.Columns {
				values_array = append(values_array, vv.FieldName+" "+vv.DataType)
			}
			CreateTable(v.Table_Name, values_array)
		}
		for _, v := range mig.Add_Column {
			for _, vv := range v.Columns {
				AddColumn(v.Table_Name, vv.FieldName, vv.DataType)
			}
		}
		for _, v := range mig.Drop_Column {
			for _, vv := range v.Columns {
				RemoveColumn(v.Table_Name, vv.FieldName)
			}
		}
		for _, v := range mig.Add_Index {
			var fieldname_array []string
			for _, vv := range v.Columns {
				fieldname_array = append(fieldname_array, vv.FieldName)
			}
			AddIndex(v.Table_Name, v.Index_Type, fieldname_array)
		}
		for _, v := range mig.Drop_Index {
			var fieldname_array []string
			for _, vv := range v.Columns {
				fieldname_array = append(fieldname_array, vv.FieldName)
			}
			RemoveIndex(v.Table_Name, v.Index_Type, fieldname_array)
		}
	}
}
func execQuery(query string) {
	fmt.Println("Postgres---" + query)
	// q, err := Db.Query(query)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer q.Close()
}

func CreateTable(table_name string, field_datatype []string) {
	query := "CREATE TABLE " + table_name + " (" + strings.Join(field_datatype, ",") + ")"
	execQuery(query)
	return
}

func DropTable(table_name string) {
	query := "DROP TABLE " + table_name
	execQuery(query)
	return
}

func AddColumn(table_name string, column_name string, data_type string) {
	query := "ALTER TABLE " + table_name + " ADD " + column_name + " " + data_type
	execQuery(query)
	return
}

func RemoveColumn(table_name string, column_name string) {
	query := "ALTER TABLE " + table_name + " DROP " + column_name
	execQuery(query)
	return
}

func AddIndex(table_name string, index_type string, field []string) {
	sort.Strings(field)
	tmp_index_name := strings.ToLower(strings.Join(field, "_") + "_index")
	query := "CREATE " + strings.ToUpper(index_type) + " INDEX " + tmp_index_name + " ON " + table_name + "( " + strings.Join(field, ",") + " )"
	execQuery(query)
	return
}

func RemoveIndex(table_name string, index_type string, field []string) {
	sort.Strings(field)
	tmp_index_name := strings.ToLower(strings.Join(field, "_") + "_index")
	query := ""
	if index_type != "" {
		query = "ALTER TABLE " + table_name + " DROP " + strings.ToUpper(index_type)
	} else {
		query = "ALTER TABLE " + table_name + " DROP INDEX " + tmp_index_name
	}
	execQuery(query)
	return
}
