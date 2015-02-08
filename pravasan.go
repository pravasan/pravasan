package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/howeyc/gopass"
)

const (
	currentVersion = "0.3"
	layout         = "20060102150405"

	// FieldDataTypeRegexp contains Regular Expression to split field name & field data type.
	FieldDataTypeRegexp = `^([A-Za-z_0-9$]{2,15}):([A-Za-z]{2,15})`
	infoText            = "\033[97m[\033[0;36mINFO\033[97m] "
	resetText           = "\033[0m"
	runningText         = "\033[97m[\033[33mRUNNING\033[97m] "
	successText         = "\033[97m[\033[32mSUCCESS\033[97m] "
	errorText           = "\033[97m[\033[31mERROR\033[97m] "
	doneText            = "\033[97m[\033[32mDONE\033[97m] "
)

// MigInterface #TODO(kishorevaishnav): need to write some comment
type MigInterface interface {
	Init(Config)
	GetLastMigrationNo() string
	CreateMigrationTable()
	ProcessNow(Migration, UpDown, string, bool)
}

var (
	flagPassword    = flag.Bool("p", false, "database password")
	version         = flag.Bool("version", false, "print Pravasan version")
	storeDirectSQL  = flag.Bool("storeDirectSQL", false, "Store SQL in migration file instead of XML / JSON")
	autoAddColumns  = flag.String("autoAddColumns", "", "Add default columns when table is created")
	configOutput    = flag.String("confOutput", currentConfFileFormat, "config file format: json, xml")
	dbHostname      = flag.String("h", "localhost", "database hostname, default: localhost")
	dbName          = flag.String("d", "", "database name")
	dbPort          = flag.String("port", "3306", "database port, default: 3306")
	dbType          = flag.String("dbType", "mysql", "database type, default: mysql")
	dbUsername      = flag.String("u", "", "database username")
	indexPrefix     = flag.String("indexPrefix", "idx", "prefix for creating Indexes, default: idx")
	indexSuffix     = flag.String("indexSuffix", "", "suffix for creating Indexes")
	migDir          = flag.String("migDir", "./", "migration file stored directory, default: ./ ")
	migFileExtn     = flag.String("migFileExtn", "prvsn", "migration file extension, default: prvsn")
	migFilePrefix   = flag.String("migFilePrefix", "", "prefix for migration file")
	migOutputFormat = flag.String("migOutput", "json", "current supported format: json, xml & deafult: json")
	migTableName    = flag.String("migTableName", "schema_migrations", "migration table name, default: schema_migrations")

	argArray              []string
	config                Config
	currentConfFileFormat = "json"
	m                     Migration
	Db                    *sql.DB
	localConfig           Config
	localUpDown           string
	migrationTableName    string
	workingVersion        string
	err                   error

	gdmMy MySQLStruct
	gdmPq PostgresStruct
	gdmSl SQLite3Struct
)

func main() {
	initializeDefaults()

	if strings.LastIndex(os.Args[0], "pravasan") < 1 || len(argArray) == 0 {
		fmt.Println(errorText + "wrong usage, or no arguments specified." + resetText)
		return
	}
	_, err := migrationDirectoryExists()
	if err != nil {
		fmt.Println(err)
		return
	}
	switch argArray[0] {
	case "add", "a":
		fn, mm := generateMigration(argArray)
		fmt.Println(fn)
		writeToFile(fn, mm, config.MigrationOutputFormat)
	case "create", "c":
		if len(argArray) > 1 && argArray[1] != "" && argArray[1] == "conf" {
			createConfigurationFile()
		} else {
			createMigration()
		}
	case "down", "d":
		migrateUpDown("down")
	case "up", "u":
		migrateUpDown("up")
	default:
		panic(errorText + "No or Wrong Actions provided." + resetText)
	}
	fmt.Println(doneText + "Completed." + resetText)
	os.Exit(1)
}

func migrationDirectoryExists() (string, error) {
	if _, err := os.Stat(config.MigrationDirectory); err == nil {
		return config.MigrationDirectory, nil
	}
	e := errors.New(config.MigrationDirectory + " doesn't exists.")
	return config.MigrationDirectory, e
}

func initializeDefaults() {
	checkConfigFileExists(&config)

	flag.Parse()

	if *version {
		fmt.Println(printCurrentVersion())
		if len(flag.Args()) == 0 {
			os.Exit(0)
		}
	}
	if *flagPassword {
		fmt.Printf("Enter DB Password : ")
		pw := gopass.GetPasswd()
		config.DbPassword = string(pw)
	}

	if *storeDirectSQL {
		config.StoreDirectSQL = "true"
	}

	config.AutoAddColumns = updateConfigValue(config.AutoAddColumns, *autoAddColumns, "")
	config.DbHostname = updateConfigValue(config.DbHostname, *dbHostname, "localhost")
	config.DbName = updateConfigValue(config.DbName, *dbName, "")
	config.DbPort = updateConfigValue(config.DbPort, *dbPort, "3306")
	config.DbType = updateConfigValue(config.DbType, *dbType, "mysql")
	config.DbUsername = updateConfigValue(config.DbUsername, *dbUsername, "")
	config.IndexPrefix = updateConfigValue(config.IndexPrefix, *indexPrefix, "idx")
	config.IndexSuffix = updateConfigValue(config.IndexSuffix, *indexSuffix, "")
	config.MigrationDirectory = updateConfigValue(config.MigrationDirectory, strings.Trim(*migDir, "/")+"/", "./")
	config.MigrationFileExtension = updateConfigValue(config.MigrationFileExtension, *migFileExtn, "prvsn")
	config.MigrationFilePrefix = updateConfigValue(config.MigrationFilePrefix, *migFilePrefix, "")
	config.MigrationOutputFormat = updateConfigValue(config.MigrationOutputFormat, strings.ToLower(*migOutputFormat), "json")
	config.MigrationTableName = updateConfigValue(config.MigrationTableName, strings.ToLower(*migTableName), "schema_migrations")
	if *dbType == "sqlite3" {
		config.DbHostname = ""
		config.DbPort = ""
		config.DbUsername = ""
		config.DbPassword = ""
	}

	argArray = flag.Args()
}

func updateConfigValue(originalValue string, overwriteValue string, defValue string) string {
	if overwriteValue != "" && overwriteValue != defValue {
		return overwriteValue
	}
	if originalValue == "" {
		if overwriteValue != "" {
			return overwriteValue
		}
		return defValue
	}
	return originalValue
}

func setDB() (mi MigInterface) {
	switch config.DbType {
	case "postgres":
		// midb := gdmPq
		mi = MigInterface(gdmPq)
	case "sqlite3":
		// midb := gdmSl{}
		mi = MigInterface(gdmSl)
	default:
		// midb := gdmMy{}
		mi = MigInterface(gdmMy)
	}
	mi.Init(config)
	return mi
}

func createMigration() {
	abc := setDB()
	abc.CreateMigrationTable()
}

func generateMigration(argsArray []string) (filename string, mm Migration) {
	t := time.Now()
	// mm = m.Migration{}
	mm.ID = t.Format(layout)
	switch argsArray[1] {
	case "add_column", "ac":
		fnAddColumn(&mm.Up, argsArray[2], argsArray[3:])
		fnDropColumn(&mm.Down, argsArray[2], argsArray[3:])
	case "add_index", "ai":
		fnAddIndex(&mm.Up, &mm.Down, argsArray[2], argsArray[3:])
	case "create_table", "ct":
		fnCreateTable(&mm.Up, argsArray[2], argsArray[3:])
		fnDropTable(&mm.Down, argsArray[2])
	case "drop_column", "dc":
		fnDropColumn(&mm.Up, argsArray[2], argsArray[3:])
		fnAddColumn(&mm.Down, argsArray[2], argsArray[3:])
	case "drop_index", "di":
		fnAddIndex(&mm.Down, &mm.Up, argsArray[2], argsArray[3:])
	case "drop_table", "dt":
		fnDropTable(&mm.Up, argsArray[2])
		fnCreateTable(&mm.Down, argsArray[2], argsArray[3:])
	case "rename_table", "rt":
		fnRenameTable(&mm.Up, &mm.Down, argsArray[2], argsArray[3])
	case "sql", "s":
		fnSql(&mm)
	// case "change_column", "cc":
	// 	fnChangeColumn(&mm.Up, &mm.Down)
	// case "rename_column", "rc":
	// 	fnRenameColumn(&mm.Up, &mm.Down)
	default:
		panic("No or wrong Actions provided.")
	}
	filename = config.MigrationDirectory + config.MigrationFilePrefix + mm.ID + "." + config.MigrationOutputFormat + "." + config.MigrationFileExtension
	return
}

func fnSql(mig *Migration) {
	fmt.Println("Hint : Type as many lines as you want, when you want to finish ^D (non Windows) or ^X (Windows)")
	var fp *os.File
	fp = os.Stdin
	reader := bufio.NewReaderSize(fp, 4096)

	fmt.Println("Enter SQL statements for Up section of migration")
	var localSql string
	for {
		line, _, err := reader.ReadLine()
		localSql = localSql + string(line) + " "
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}
	mig.Up.Sql = strings.TrimSpace(localSql)

	fmt.Println("Now enter SQL statements for Down section of migration")
	localSql = ""
	for {
		line, _, err := reader.ReadLine()
		localSql = localSql + string(line) + " "
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}
	mig.Down.Sql = strings.TrimSpace(localSql)
}

func fnAddColumn(mm *UpDown, tableName string, fieldArray []string) {
	ac := AddColumn{TableName: tableName}
	for key, value := range fieldArray {
		fieldArray[key] = strings.Trim(value, ", ")
		if r, _ := regexp.Compile(FieldDataTypeRegexp); r.MatchString(fieldArray[key]) == true {
			split := r.FindAllStringSubmatch(fieldArray[key], -1)
			col := Columns{
				FieldName: split[0][1],
				DataType:  split[0][2]}
			ac.Columns = append(ac.Columns, col)
		} else {
			ac = AddColumn{}
		}
	}
	if config.StoreDirectSQL == "true" {
		a := UpDown{}
		a.AddColumn = append(a.AddColumn, ac)
		mm.Sql = gdmMy.ReturnQuery(a)
	} else {
		mm.AddColumn = append(mm.AddColumn, ac)
	}
}

func fnAddIndex(mUp *UpDown, mDown *UpDown, tableName string, fieldArray []string) {
	ai := AddIndex{TableName: tableName}
	di := DropIndex{TableName: tableName}
	for key, value := range fieldArray {
		fieldArray[key] = strings.Trim(value, ", ")
		col := Columns{}
		if r, _ := regexp.Compile(FieldDataTypeRegexp); r.MatchString(fieldArray[key]) == true {
			split := r.FindAllStringSubmatch(fieldArray[key], -1)
			col.FieldName = split[0][1]
			col.DataType = split[0][2]
			ai.Columns = append(ai.Columns, col)
			di.Columns = append(di.Columns, col)
		} else if fieldArray[key] != "" {
			col.FieldName = fieldArray[key]
			ai.Columns = append(ai.Columns, col)
			di.Columns = append(di.Columns, col)
		} else {
			ai = AddIndex{}
			di = DropIndex{}
		}
	}
	if config.StoreDirectSQL == "true" {
		a := UpDown{}
		a.AddIndex = append(a.AddIndex, ai)
		mUp.Sql = gdmMy.ReturnQuery(a)
		b := UpDown{}
		b.DropIndex = append(b.DropIndex, di)
		mDown.Sql = gdmMy.ReturnQuery(a)
	} else {
		mUp.AddIndex = append(mUp.AddIndex, ai)
		mDown.DropIndex = append(mDown.DropIndex, di)
	}
}

func fnChangeColumn(mUp *UpDown, mDown *UpDown) {
}

func fnCreateTable(mm *UpDown, tableName string, fieldArray []string) {
	ct := CreateTable{TableName: tableName}
	providedCol := []Columns{}
	for key, value := range fieldArray {
		fieldArray[key] = strings.Trim(value, ", ")
		if r, _ := regexp.Compile(FieldDataTypeRegexp); r.MatchString(fieldArray[key]) == true {
			split := r.FindAllStringSubmatch(fieldArray[key], -1)
			col := Columns{
				FieldName: split[0][1],
				DataType:  split[0][2]}
			providedCol = append(providedCol, col)
		}
	}

	fmt.Println("+++++ ", config.AutoAddColumns)
	// Below snippet will add extra columns automatically.
	if config.AutoAddColumns != "" {
		aacArray := strings.Fields(config.AutoAddColumns)
		for _, value := range aacArray {
			flag := true
			var autoFieldName, autoDataType string
			if r, _ := regexp.Compile(FieldDataTypeRegexp); r.MatchString(value) == true {
				split := r.FindAllStringSubmatch(value, -1)
				autoFieldName, autoDataType = split[0][1], split[0][2]
				for _, pcValue := range providedCol {
					if pcValue.FieldName == split[0][1] {
						flag = false
						break
					}
				}
			} else {
				flag = false
			}
			if flag == true {
				ct.Columns = append(ct.Columns, Columns{FieldName: autoFieldName, DataType: autoDataType})
			}
		}
	}
	ct.Columns = append(ct.Columns, providedCol...)
	if config.StoreDirectSQL == "true" {
		a := UpDown{}
		a.CreateTable = append(a.CreateTable, ct)
		mm.Sql = gdmMy.ReturnQuery(a)
	} else {
		mm.CreateTable = append(mm.CreateTable, ct)
	}
}

func fnDropColumn(mm *UpDown, tableName string, fieldArray []string) {
	dc := DropColumn{TableName: tableName}
	for key, value := range fieldArray {
		fieldArray[key] = strings.Trim(value, ", ")
		col := Columns{}
		if r, _ := regexp.Compile(FieldDataTypeRegexp); r.MatchString(fieldArray[key]) == true {
			split := r.FindAllStringSubmatch(fieldArray[key], -1)
			col.FieldName = split[0][1]
			col.DataType = split[0][2]
			dc.Columns = append(dc.Columns, col)
		} else if fieldArray[key] != "" {
			col.FieldName = fieldArray[key]
			dc.Columns = append(dc.Columns, col)
		}
	}
	if config.StoreDirectSQL == "true" {
		a := UpDown{}
		a.DropColumn = append(a.DropColumn, dc)
		mm.Sql = gdmMy.ReturnQuery(a)
	} else {
		mm.DropColumn = append(mm.DropColumn, dc)
	}
}

func fnDropTable(mm *UpDown, tableName string) {
	dt := DropTable{TableName: tableName}
	if config.StoreDirectSQL == "true" {
		a := UpDown{}
		a.DropTable = append(a.DropTable, dt)
		mm.Sql = gdmMy.ReturnQuery(a)
	} else {
		mm.DropTable = append(mm.DropTable, dt)
	}
}

func fnRenameTable(mUp *UpDown, mDown *UpDown, srcTableName string, destTablename string) {
	mUp.RenameTable = append(mUp.RenameTable, RenameTable{
		OldTableName: srcTableName,
		NewTableName: destTablename})
	mDown.RenameTable = append(mDown.RenameTable, RenameTable{
		OldTableName: destTablename,
		NewTableName: srcTableName})
}

// migrateUpDown - used to perform the Action Migration either Up or Down & chooses the DB too.
func migrateUpDown(updown string) {

	// Actual migration can't happen without having DB Name & DB User, rest can be default.
	if config.DbName == "" || (config.DbUsername == "" && *dbType != "sqlite3") {
		fmt.Println(errorText + "Either Database Name or Username is not mentioned, or both are missed to mention." + resetText)
		return
	}

	var (
		err          error
		files        []string
		processCount int
		reverseCount = 1
		force        = false
	)

	fmt.Println(infoText + "Collecting migration files..." + resetText)
	files = migrationFiles(updown)
	if 0 == len(files) {
		fmt.Println("No migration files present.")
		return
	}

	if len(argArray) > 1 && argArray[1] != "" {
		if updown == "down" {
			reverseCount, err = strconv.Atoi(strings.Replace(argArray[1], "-", "", -1))
			if err != nil {
				log.Println(errorText + "Wrong count to be reversed, only integer values accepted." + resetText)
				log.Fatal(err)
			}
		} else if updown == "up" {
			files, err = checkMigrationFilesExists(argArray[1:len(argArray)])
			if err != nil {
				fmt.Println(errorText + "Any one of the version number of the file mentioned is wrong." + resetText)
				return
			}
			force = true
		}
	}

	for _, filename := range files {
		// During migration down if count is reached then exit
		if updown == "down" && reverseCount == processCount {
			break
		}

		// Read the content of the file & store into the mm structure
		fmt.Print(runningText + filename + resetText + " - ")
		bs, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Println(err)
			return
		}
		docScript := []byte(bs)

		var (
			mm  Migration
			mig UpDown
		)

		if config.MigrationOutputFormat == "json" {
			json.Unmarshal(docScript, &mm)
		} else {
			xml.Unmarshal(docScript, &mm)
		}

		if updown == "up" {
			mig = mm.Up
		} else {
			mig = mm.Down
		}

		if config.DbType == "mysql" {
			abc := setDB()
			if updown == "down" && mm.ID > abc.GetLastMigrationNo() {
				continue
			}
			// gdmMy.ProcessNow(mm, mig, updown, force)
			abc.ProcessNow(mm, mig, updown, force)
		} else {
			abc := setDB()
			if updown == "down" && mm.ID > abc.GetLastMigrationNo() {
				continue
			}
			abc.ProcessNow(mm, mig, updown, force)
		}
		processCount++
		fmt.Println(successText + filename + " Migrated." + resetText)
	}
}

func checkMigrationFilesExists(verFiles []string) (orderVerFiles []string, err error) {
	err = nil
	for _, indvlFile := range verFiles {
		indvlFile = strings.Replace(indvlFile, ",", "", -1)
		indvlFile = strings.Replace(indvlFile, " ", "", -1)
		if !strings.HasSuffix(indvlFile, "."+config.MigrationOutputFormat+"."+config.MigrationFileExtension) {
			indvlFile = indvlFile + "." + config.MigrationOutputFormat + "." + config.MigrationFileExtension
		}
		if config.MigrationFilePrefix != "" && !strings.HasPrefix(indvlFile, config.MigrationFilePrefix) {
			indvlFile = config.MigrationFilePrefix + indvlFile
		}
		if _, err := os.Stat(indvlFile); err == nil {
			orderVerFiles = append(orderVerFiles, indvlFile)
		} else {
			err = errors.New(errorText + "coudn't able to read file: " + indvlFile + resetText)
			break
		}
	}
	if err != nil {
		return nil, err
	}
	sort.Strings(orderVerFiles)
	return orderVerFiles, err
}

func migrationFiles(updown string) []string {
	files, _ := ioutil.ReadDir(config.MigrationDirectory)
	var onlyMigFiles []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), "."+config.MigrationOutputFormat+"."+config.MigrationFileExtension) && strings.HasPrefix(f.Name(), config.MigrationFilePrefix) {
			onlyMigFiles = append(onlyMigFiles, config.MigrationDirectory+f.Name())
		}
	}

	if updown == "down" {
		sort.Sort(sort.Reverse(sort.StringSlice(onlyMigFiles)))
	} else {
		sort.Strings(onlyMigFiles)
	}
	return onlyMigFiles
}

func writeToFile(filename string, obj interface{}, format string) {
	var content []byte
	if format == "json" {
		// Indenting the JSON format
		content, _ = json.MarshalIndent(obj, " ", "  ")
	} else if format == "xml" {
		// Indenting the XML format
		content, _ = xml.MarshalIndent(obj, " ", "  ")
	} else {
		return
	}

	// Write to a new File.
	file, _ := os.Create(filename)
	file.Write(content)
	file.Close()
}

func createConfigurationFile() {
	writeToFile("pravasan.conf."+currentConfFileFormat, config, currentConfFileFormat)
	fmt.Println(successText + "Config file created / updated.")
}

func printCurrentVersion() string {
	return infoText + "pravasan version " + currentVersion + resetText
}

// checkConfigFileExists loads the configuration if it exists whether it is XML or JSON.
func checkConfigFileExists(c *Config) {
	if _, err := os.Stat("./pravasan.conf.json"); err == nil {
		currentConfFileFormat = "json"
		bs, err := ioutil.ReadFile("pravasan.conf.json")
		if err != nil {
			fmt.Println(err)
			return
		}
		docScript := []byte(bs)
		json.Unmarshal(docScript, &c)
	} else if _, err := os.Stat("./pravasan.conf.xml"); err == nil {
		currentConfFileFormat = "xml"
		bs, err := ioutil.ReadFile("pravasan.conf.xml")
		if err != nil {
			fmt.Println(err)
			return
		}
		docScript := []byte(bs)
		xml.Unmarshal(docScript, &c)
	}
}
