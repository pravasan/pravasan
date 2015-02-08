package main

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	// #TODO(kishorevaishnav): need to write comment why use "_"
	_ "github.com/go-sql-driver/mysql"
)

//MySQLStruct #TODO(kishorevaishnav): need to write some comment & need to write different name instead of MySQLStruct
type MySQLStruct struct {
	bTQ string
}

// Init is called to initiate the connection to check and do some activities
func (s MySQLStruct) Init(c Config) {
	// This can be useful to check for version and any other dependencies etc.,
	// fmt.Println("mysql init() it runs before other functions")
	if c.DbPort == "" {
		c.DbPort = "3306"
	}
	Db, _ = sql.Open("mysql", c.DbUsername+":"+c.DbPassword+"@/"+c.DbName)
	migrationTableName = c.MigrationTableName

	localConfig = c
	s.bTQ = "`" // s.bTQ = backTickQuote
}

// GetLastMigrationNo - get the last migration it has executed.
func (s MySQLStruct) GetLastMigrationNo() string {
	maxVersion := ""
	query := "SELECT max(" + s.bTQ + "version" + s.bTQ + ") FROM " + s.bTQ + migrationTableName + s.bTQ
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
func (s MySQLStruct) CreateMigrationTable() {
	query := "CREATE TABLE " + s.bTQ + migrationTableName + s.bTQ + " (" + s.bTQ + "version" + s.bTQ + " VARCHAR(15))"
	q, err := Db.Query(query)
	defer q.Close()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Table Created Successfully.")
	}
}

// ProcessNow is used to run the actual migraition whether it is UP or DOWN.
func (s MySQLStruct) ProcessNow(lm Migration, mig UpDown, updown string, force bool) {
	if updown == "up" {
		if force == false && lm.ID <= s.GetLastMigrationNo() {
			return
		}
		if force == true && s.checkMigrationExecutedForID(lm.ID) {
			fmt.Println(lm.ID + " -> Its already executed.")
			return
		}
	}

	localUpDown = updown
	workingVersion = lm.ID

	if nid, _ := strconv.Atoi(lm.ID); nid != 0 {
		fmt.Println("Executing ID : ", lm.ID)
		s.execQuery(s.ReturnQuery(mig))
		if mig.Sql != "" {
			s.directSQL(mig.Sql)
		}
		s.updateMigrationTable()
	}
}

// ReturnQuery will return direct SQL query
func (s MySQLStruct) ReturnQuery(mig UpDown) string {
	for _, v := range mig.AddColumn {
		for _, vv := range v.Columns {
			// #TODO(kishorevaishnav): need to remove the return out of the for loop
			return s.addColumn(s.bTQ+v.TableName+s.bTQ, s.bTQ+vv.FieldName+s.bTQ+" ", s.dataTypeConversion(vv.DataType))
		}
	}
	for _, v := range mig.AddIndex {
		var fieldNameArray []string
		for _, vv := range v.Columns {
			fieldNameArray = append(fieldNameArray, s.bTQ+vv.FieldName+s.bTQ+" ")
		}
		return s.addIndex(s.bTQ+v.TableName+s.bTQ, v.IndexType, fieldNameArray)
	}
	for _, v := range mig.CreateTable {
		var valuesArray []string
		for _, vv := range v.Columns {
			valuesArray = append(valuesArray, s.bTQ+vv.FieldName+s.bTQ+" "+s.dataTypeConversion(vv.DataType))
		}
		return s.createTable(s.bTQ+v.TableName+s.bTQ, valuesArray)
	}
	for _, v := range mig.DropColumn {
		for _, vv := range v.Columns {
			// #TODO(kishorevaishnav): need to remove the return out of the for loop
			return s.dropColumn(s.bTQ+v.TableName+s.bTQ, s.bTQ+vv.FieldName+s.bTQ+" ")
		}
	}
	for _, v := range mig.DropIndex {
		var fieldNameArray []string
		for _, vv := range v.Columns {
			fieldNameArray = append(fieldNameArray, s.bTQ+vv.FieldName+s.bTQ+" ")
		}
		return s.dropIndex(s.bTQ+v.TableName+s.bTQ, v.IndexType, fieldNameArray)
	}
	for _, v := range mig.DropTable {
		// #TODO(kishorevaishnav): need to remove the return out of the for loop
		return s.dropTable(s.bTQ + v.TableName + s.bTQ)
	}
	for _, v := range mig.RenameTable {
		// #TODO(kishorevaishnav): need to remove the return out of the for loop
		return s.renameTable(s.bTQ+v.OldTableName+s.bTQ, s.bTQ+v.NewTableName+s.bTQ)
	}
	return ""
}

func (s MySQLStruct) updateMigrationTable() {
	var query string
	if localUpDown == "up" {
		query = "INSERT INTO " + s.bTQ + migrationTableName + s.bTQ + "(" + s.bTQ + "version" + s.bTQ + ") VALUES ('" + workingVersion + "')"
	} else {
		query = "DELETE FROM " + s.bTQ + migrationTableName + s.bTQ + " WHERE " + s.bTQ + "version" + s.bTQ + "='" + workingVersion + "'"
	}
	q, err := Db.Query(query)
	defer q.Close()
	if err != nil {
		log.Println("not able to add version to the existing migration table")
		log.Fatal(err)
	}
}

func (s MySQLStruct) checkMigrationExecutedForID(id string) bool {
	var version = ""
	query := "SELECT " + s.bTQ + "version" + s.bTQ + " FROM " + s.bTQ + migrationTableName + s.bTQ + " WHERE " + s.bTQ + "version" + s.bTQ + "='" + id + "'"
	q, err := Db.Query(query)
	defer q.Close()
	if err != nil {
		log.Println("couldn't able to execute the check version query...")
		log.Fatal(err)
	} else {
		q.Next()
		q.Scan(&version)
	}
	if version == "" {
		return false
	}
	return true
}

func (s MySQLStruct) dataTypeConversion(dt string) string {
	switch dt {
	case "string":
		return "VARCHAR(255)"
	case "int":
		return "INTEGER"
	}
	return dt
}

func (s MySQLStruct) directSQL(query string) {
	s.execQuery(query)
	return
}

func (s MySQLStruct) execQuery(query string) {
	fmt.Println("MySQL---" + query)
	q, err := Db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer q.Close()
}

func (s MySQLStruct) createTable(tableName string, fieldDataType []string) string {
	return "CREATE TABLE " + tableName + " (" + strings.Join(fieldDataType, ",") + ")"
}

func (s MySQLStruct) dropTable(tableName string) string {
	return "DROP TABLE " + tableName
}

func (s MySQLStruct) addColumn(tableName string, columnName string, dataType string) string {
	return "ALTER TABLE " + tableName + " ADD " + columnName + " " + dataType
}

func (s MySQLStruct) dropColumn(tableName string, columnName string) string {
	return "ALTER TABLE " + tableName + " DROP " + columnName
}

func (s MySQLStruct) addIndex(tableName string, indexType string, field []string) string {
	// #TODO(kishorevaishnav): currently indexType is always empty as we don't have a proper way.

	sort.Strings(field)
	tmpIndexName := localConfig.IndexPrefix + "_" + strings.Join(field, "_") + "_" + localConfig.IndexSuffix
	tmpIndexName = strings.Trim(strings.Replace(strings.Replace(strings.ToLower(tmpIndexName), s.bTQ+"", "", -1), " ", "", -1), "_")
	return "CREATE " + strings.ToUpper(indexType) + " INDEX " + tmpIndexName + " ON " + tableName + "( " + strings.Join(field, ",") + " )"
}

func (s MySQLStruct) dropIndex(tableName string, indexType string, field []string) string {
	// #TODO(kishorevaishnav): currently indexType is always empty as we don't have a proper way.

	sort.Strings(field)
	tmpIndexName := localConfig.IndexPrefix + "_" + strings.Join(field, "_") + "_" + localConfig.IndexSuffix
	tmpIndexName = strings.Trim(strings.Replace(strings.Replace(strings.ToLower(tmpIndexName), s.bTQ+"", "", -1), " ", "", -1), "_")
	if indexType != "" {
		return "ALTER TABLE " + tableName + " DROP " + strings.ToUpper(indexType)
	}
	return "ALTER TABLE " + tableName + " DROP INDEX " + tmpIndexName
}

func (s MySQLStruct) renameTable(oldTableName string, newTableName string) string {
	return "ALTER TABLE " + oldTableName + " RENAME " + newTableName
}
