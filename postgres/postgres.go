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

// All global variable declartion done here.
var (
	Db                 *sql.DB
	workingVersion     string
	localUpDown        string
	migrationTableName string
)

func init() {
	// This can be useful to check for version and any other dependencies etc.,
	// fmt.Println("mysql init() it runs before other functions")
}

// Init is called to initiate the connection to check and do some activities
func Init(c m.Config) {
	var err error
	Db, err = sql.Open("postgres", "postgres://"+c.DbUsername+":"+c.DbPassword+"@"+c.DbHostname+"/"+c.DbName+"?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	migrationTableName = c.MigrationTableName
}

// GetLastMigrationNo to get what is the last migration it has executed.
func GetLastMigrationNo() string {
	var maxVersion = ""
	query := "SELECT max(\"version\") FROM \"" + migrationTableName + "\""
	q, err := Db.Query(query)
	defer q.Close()
	if err != nil {
		log.Println(migrationTableName + " table doesn't exists")
		log.Fatal(err)
	} else {
		q.Next()
		q.Scan(&maxVersion)
	}
	return maxVersion
}

// CreateMigrationTable used to create the schema_migration if it doesn't exists.
func CreateMigrationTable() {
	query := "CREATE TABLE \"" + migrationTableName + "\" (version VARCHAR(255))"
	q, err := Db.Query(query)
	defer q.Close()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Table Created Successfully.")
	}
}

func updateMigrationTable() {
	var query string
	if localUpDown == "up" {
		query = "INSERT INTO \"" + migrationTableName + "\"(version) VALUES ('" + workingVersion + "')"
	} else {
		query = "DELETE FROM \"" + migrationTableName + "\" WHERE version='" + workingVersion + "'"
	}
	q, err := Db.Query(query)
	defer q.Close()
	if err != nil {
		log.Println("not able to add version to the existing migration table")
		log.Fatal(err)
	}
}

func dataTypeConversion(dt string) string {
	switch dt {
	case "string":
		return "VARCHAR(255)"
	case "int":
		return "INTEGER"
	}
	return dt
}

// ProcessNow is used to run the actual migraition whether it is UP or DOWN.
func ProcessNow(lm m.Migration, mig m.UpDown, updown string) {
	if updown == "up" && lm.ID <= GetLastMigrationNo() {
		return
	}
	localUpDown = updown

	workingVersion = lm.ID
	nid, _ := strconv.Atoi(lm.ID)
	if nid != 0 {
		fmt.Println("Executing ID : ", lm.ID)
		for _, v := range mig.CreateTable {
			var valuesArray []string
			for _, vv := range v.Columns {
				valuesArray = append(valuesArray, "\""+vv.FieldName+"\" "+dataTypeConversion(vv.DataType))
			}
			createTable("\""+v.TableName+"\"", valuesArray)
		}
		for _, v := range mig.DropTable {
			dropTable("\"" + v.TableName + "\"")
		}
		for _, v := range mig.AddColumn {
			for _, vv := range v.Columns {
				addColumn("\""+v.TableName+"\"", "\""+vv.FieldName+"\" ", dataTypeConversion(vv.DataType))
			}
		}
		for _, v := range mig.DropColumn {
			for _, vv := range v.Columns {
				dropColumn("\""+v.TableName+"\"", "\""+vv.FieldName+"\" ")
			}
		}
		for _, v := range mig.AddIndex {
			var fieldNameArray []string
			for _, vv := range v.Columns {
				fieldNameArray = append(fieldNameArray, "\""+vv.FieldName+"\" ")
			}
			addIndex("\""+v.TableName+"\"", v.IndexType, fieldNameArray)
		}
		for _, v := range mig.DropIndex {
			var fieldNameArray []string
			for _, vv := range v.Columns {
				fieldNameArray = append(fieldNameArray, "\""+vv.FieldName+"\" ")
			}
			dropIndex("\""+v.TableName+"\"", v.IndexType, fieldNameArray)
		}
		if mig.Sql != "" {
			directSQL(mig.Sql)
		}
		updateMigrationTable()
	}
}

func directSQL(query string) {
	execQuery(query)
	return
}

func execQuery(query string) {
	fmt.Println("Postgres---" + query)
	q, err := Db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer q.Close()
}

func createTable(tableName string, fieldDataType []string) {
	query := "CREATE TABLE " + tableName + " (" + strings.Join(fieldDataType, ",") + ")"
	execQuery(query)
	return
}

func dropTable(tableName string) {
	query := "DROP TABLE " + tableName
	execQuery(query)
	return
}

func addColumn(tableName string, columnName string, dataType string) {
	query := "ALTER TABLE " + tableName + " ADD " + columnName + " " + dataType
	execQuery(query)
	return
}

func dropColumn(tableName string, columnName string) {
	query := "ALTER TABLE " + tableName + " DROP " + columnName
	execQuery(query)
	return
}

func addIndex(tableName string, indexType string, field []string) {
	sort.Strings(field)
	tmpIndexName := strings.ToLower(strings.Join(field, "_") + "_index")
	query := "CREATE " + strings.ToUpper(indexType) + " INDEX " + tmpIndexName + " ON " + tableName + "( " + strings.Join(field, ",") + " )"
	execQuery(query)
	return
}

func dropIndex(tableName string, indexType string, field []string) {
	sort.Strings(field)
	tmpIndexName := strings.ToLower(strings.Join(field, "_") + "_index")
	query := ""
	if indexType != "" {
		query = "ALTER TABLE " + tableName + " DROP " + strings.ToUpper(indexType)
	} else {
		query = "ALTER TABLE " + tableName + " DROP INDEX " + tmpIndexName
	}
	execQuery(query)
	return
}
