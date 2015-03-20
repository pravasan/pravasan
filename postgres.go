package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	// #TODO(kishorevaishnav): need to write comment why use "_"
	_ "github.com/lib/pq"
)

//PostgresStruct #TODO(kishorevaishnav): need to write some comment & need to write different name instead of PostgresStruct
type PostgresStruct struct {
	bTQ string
}

//PostgreSQLSupportedDataTypes #TODO(kishorevaishnav): need to write some comment
var PostgreSQLSupportedDataTypes = map[string]string{
	"BIGINT":        "BIGINT",
	"BIGSERIAL":     "BIGSERIAL",
	"BIT":           "BIT",
	"BOOLEAN":       "BOOLEAN",
	"BOX":           "BOX",
	"BYTEA":         "BYTEA",
	"CHARACTER":     "CHARACTER",
	"CIDR":          "CIDR",
	"CIRCLE":        "CIRCLE",
	"DATE":          "DATE",
	"INET":          "INET",
	"INTEGER":       "INTEGER",
	"JSON":          "JSON",
	"LINE":          "LINE",
	"LSEG":          "LSEG",
	"MACADDR":       "MACADDR",
	"MONEY":         "MONEY",
	"NUMERIC":       "NUMERIC",
	"PATH":          "PATH",
	"POINT":         "POINT",
	"POLYGON":       "POLYGON",
	"REAL":          "REAL",
	"SMALLINT":      "SMALLINT",
	"SMALLSERIAL":   "SMALLSERIAL",
	"SERIAL":        "SERIAL",
	"TEXT":          "TEXT",
	"TSQUERY":       "TSQUERY",
	"TSVECTOR":      "TSVECTOR",
	"TXID_SNAPSHOT": "TXID_SNAPSHOT",
	"UUID":          "UUID",
	"XML":           "XML",
	"VARBIT":        "VARBIT",
	"BOOL":          "BOOL",
	"CHAR":          "CHAR",
	"VARCHAR":       "VARCHAR",
	"INT":           "INT",
	"INT8":          "INT8",
	"DECIMAL":       "DECIMAL",
	"TIMETZ":        "TIMETZ",
	"TIMESTAMPTZ":   "TIMESTAMPTZ",
	"DEC":           "DECIMAL",      // Alias of DECIMAL
	"FIXED":         "DECIMAL",      // Alias of DECIMAL
	"STRING":        "VARCHAR(255)", // Alias of VARCHAR
}

func init() {
	ListSuppDataTypes["PostgreSQL"] = PostgreSQLSupportedDataTypes
	// ListSuppDataTypes["Postgres"] = PostgreSQLSupportedDataTypes
}

// // ListOfSupportedDataTypes return the supported List of DataTypes.
// func (s PostgresStruct) ListOfSupportedDataTypes() (sdt map[string]string) {
// 	return PostgreSQLSupportedDataTypes
// }

// Init is called to initiate the connection to check and do some activities
func (s PostgresStruct) Init(c Config) {
	// This can be useful to check for version and any other dependencies etc.,
	// fmt.Println("postgres init() it runs before other functions")
	if c.DbPort == "" {
		c.DbPort = "5432"
	}
	Db, err = sql.Open("postgres", "postgres://"+c.DbUsername+":"+c.DbPassword+"@"+c.DbHostname+"/"+c.DbName+"?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	migrationTableName = c.MigrationTableName

	localConfig = c
	s.bTQ = "\"" // s.bTQ = backTickQuote
}

// GetLastMigrationNo to get what is the last migration it has executed.
func (s PostgresStruct) GetLastMigrationNo() string {
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
func (s PostgresStruct) CreateMigrationTable() {
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
func (s PostgresStruct) ProcessNow(lm Migration, mig UpDown, updown string, force bool) {
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
func (s PostgresStruct) ReturnQuery(mig UpDown) string {
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

func (s PostgresStruct) updateMigrationTable() {
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

func (s PostgresStruct) checkMigrationExecutedForID(id string) bool {
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

func (s PostgresStruct) dataTypeConversion(dt string) string {
	if PostgreSQLSupportedDataTypes[strings.ToUpper(dt)] == "" {
		fmt.Println("UnSupported DataType")
		os.Exit(1)
	}
	return PostgreSQLSupportedDataTypes[strings.ToUpper(dt)]
}

func (s PostgresStruct) directSQL(query string) {
	s.execQuery(query)
	return
}

func (s PostgresStruct) execQuery(query string) {
	fmt.Println("Postgres---" + query)
	q, err := Db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer q.Close()
}

func (s PostgresStruct) createTable(tableName string, fieldDataType []string) string {
	return "CREATE TABLE " + tableName + " (" + strings.Join(fieldDataType, ",") + ")"
}

func (s PostgresStruct) dropTable(tableName string) string {
	return "DROP TABLE " + tableName
}

func (s PostgresStruct) addColumn(tableName string, columnName string, dataType string) string {
	return "ALTER TABLE " + tableName + " ADD " + columnName + " " + dataType
}

func (s PostgresStruct) dropColumn(tableName string, columnName string) string {
	return "ALTER TABLE " + tableName + " DROP " + columnName
}

func (s PostgresStruct) addIndex(tableName string, indexType string, field []string) string {
	// #TODO(kishorevaishnav): currently indexType is always empty as we don't have a proper way.

	sort.Strings(field)
	tmpIndexName := localConfig.IndexPrefix + "_" + strings.Join(field, "_") + "_" + localConfig.IndexSuffix
	tmpIndexName = strings.Trim(strings.Replace(strings.Replace(strings.ToLower(tmpIndexName), s.bTQ+"", "", -1), " ", "", -1), "_")
	return "CREATE " + strings.ToUpper(indexType) + " INDEX " + tmpIndexName + " ON " + tableName + "( " + strings.Join(field, ",") + " )"
}

func (s PostgresStruct) dropIndex(tableName string, indexType string, field []string) string {
	// #TODO(kishorevaishnav): currently indexType is always empty as we don't have a proper way.

	sort.Strings(field)
	tmpIndexName := localConfig.IndexPrefix + "_" + strings.Join(field, "_") + "_" + localConfig.IndexSuffix
	tmpIndexName = strings.Trim(strings.Replace(strings.Replace(strings.ToLower(tmpIndexName), s.bTQ+"", "", -1), " ", "", -1), "_")
	if indexType != "" {
		return "ALTER TABLE " + tableName + " DROP " + strings.ToUpper(indexType)
	}
	return "ALTER TABLE " + tableName + " DROP INDEX " + tmpIndexName
}

func (s PostgresStruct) renameTable(oldTableName string, newTableName string) string {
	return "ALTER TABLE " + oldTableName + " RENAME " + newTableName
}
