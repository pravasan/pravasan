package main

import (
	"bufio"
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

	m "github.com/pravasan/pravasan/migration"

	gdm_my "github.com/pravasan/pravasan/mysql"
	gdm_pq "github.com/pravasan/pravasan/postgres"
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

var (
	argArray              []string
	config                m.Config
	currentConfFileFormat = "json"
)

func main() {
	one_init()
	if strings.LastIndex(os.Args[0], "pravasan") < 1 || len(argArray) == 0 {
		fmt.Println(errorText + "wrong usage, or no arguments specified." + resetText)
		os.Exit(1)
	}
	switch argArray[0] {
	case "add", "a":
		fn, mm := generateMigration(argArray)
		fmt.Println(fn)
		writeToFile(fn, mm, config.MigrationOutputFormat)
		// os.Exit(0)

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

func one_init() {
	// fmt.Println("pravasan init() it runs before other functions")

	checkConfigFileExists(&config)

	var (
		flagPassword = false
		version      = false

		configOutput    string
		dbHostname      string
		dbName          string
		dbPort          string
		dbType          string
		dbUsername      string
		indexPrefix     string
		indexSuffix     string
		migDir          string
		migFileExtn     string
		migFilePrefix   string
		migOutputFormat string
		migTableName    string
	)

	flag.BoolVar(&flagPassword, "p", false, "database password")
	flag.BoolVar(&version, "version", false, "print Pravasan version")
	flag.StringVar(&configOutput, "confOutput", currentConfFileFormat, "config file format: json, xml")
	flag.StringVar(&dbHostname, "h", "localhost", "database hostname, default: localhost")
	flag.StringVar(&dbName, "d", "", "database name")
	flag.StringVar(&dbPort, "port", "3306", "database port, default: 3306")
	flag.StringVar(&dbType, "dbType", "mysql", "database type, default: mysql")
	flag.StringVar(&dbUsername, "u", "", "database username")
	flag.StringVar(&indexPrefix, "indexPrefix", "idx", "prefix for creating Indexes, default: idx")
	flag.StringVar(&indexSuffix, "indexSuffix", "", "suffix for creating Indexes")
	flag.StringVar(&migDir, "migDir", "./", "migration file stored directory, default: ./ ")
	flag.StringVar(&migFileExtn, "migFileExtn", "prvsn", "migration file extension, default: prvsn")
	flag.StringVar(&migFilePrefix, "migFilePrefix", "", "prefix for migration file")
	flag.StringVar(&migOutputFormat, "migOutput", "json", "current supported format: json, xml & deafult: json")
	flag.StringVar(&migTableName, "migTableName", "schema_migrations", "migration table name, default: schema_migrations")
	flag.Parse()

	if version {
		fmt.Println(printCurrentVersion())
		if len(flag.Args()) == 0 {
			os.Exit(0)
		}
	}
	if flagPassword {
		fmt.Printf("Enter DB Password : ")
		pw := gopass.GetPasswd()
		config.DbPassword = string(pw)
	}

	config.DbHostname = updateConfigValue(config.DbHostname, dbHostname, "localhost")
	config.DbName = updateConfigValue(config.DbName, dbName, "")
	config.DbPort = updateConfigValue(config.DbPort, dbPort, "3306")
	config.DbType = updateConfigValue(config.DbType, dbType, "mysql")
	config.DbUsername = updateConfigValue(config.DbUsername, dbUsername, "")
	config.IndexPrefix = updateConfigValue(config.IndexPrefix, indexPrefix, "idx")
	config.IndexSuffix = updateConfigValue(config.IndexSuffix, indexSuffix, "")
	config.MigrationDirectory = updateConfigValue(config.MigrationDirectory, strings.Trim(migDir, "/")+"/", "./")
	config.MigrationFileExtension = updateConfigValue(config.MigrationFileExtension, migFileExtn, "prvsn")
	config.MigrationFilePrefix = updateConfigValue(config.MigrationFilePrefix, migFilePrefix, "")
	config.MigrationOutputFormat = updateConfigValue(config.MigrationOutputFormat, strings.ToLower(migOutputFormat), "json")
	config.MigrationTableName = updateConfigValue(config.MigrationTableName, strings.ToLower(migTableName), "schema_migrations")

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

func createMigration() {
	if config.DbType == "mysql" {
		gdm_my.Init(config)
		gdm_my.CreateMigrationTable()
	} else {
		gdm_pq.Init(config)
		gdm_pq.CreateMigrationTable()
	}
}

func generateMigration(ab []string) (filename string, mm m.Migration) {
	t := time.Now()
	mm = m.Migration{}
	mm.ID = t.Format(layout)
	switch ab[1] {
	case "add_column", "ac":
		fnAddColumn(&mm.Up, ab[2], ab[3:])
		fnDropColumn(&mm.Down, ab[2], ab[3:])
	case "add_index", "ai":
		fnAddIndex(&mm.Up, &mm.Down, ab[2], ab[3:])
	case "create_table", "ct":
		fnCreateTable(&mm.Up, ab[2], ab[3:])
		fnDropTable(&mm.Down, ab[2])
	case "drop_column", "dc":
		fnDropColumn(&mm.Up, ab[2], ab[3:])
		fnAddColumn(&mm.Down, ab[2], ab[3:])
	case "drop_index", "di":
		fnAddIndex(&mm.Down, &mm.Up, ab[2], ab[3:])
	case "drop_table", "dt":
		fnDropTable(&mm.Up, ab[2])
		fnCreateTable(&mm.Down, ab[2], ab[3:])
	case "rename_table", "rt":
		fnRenameTable(&mm.Up, &mm.Down, ab[2], ab[3])
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

func fnSql(mig *m.Migration) {
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

func fnAddColumn(mm *m.UpDown, tableName string, fieldArray []string) {
	ac := m.AddColumn{TableName: tableName}
	for key, value := range fieldArray {
		fieldArray[key] = strings.Trim(value, ", ")
		r, _ := regexp.Compile(FieldDataTypeRegexp)
		if r.MatchString(fieldArray[key]) == true {
			split := r.FindAllStringSubmatch(fieldArray[key], -1)
			col := m.Columns{
				FieldName: split[0][1],
				DataType:  split[0][2]}
			ac.Columns = append(ac.Columns, col)
		} else {
			ac = m.AddColumn{}
		}
	}
	mm.AddColumn = append(mm.AddColumn, ac)
}

func fnAddIndex(mUp *m.UpDown, mDown *m.UpDown, tableName string, fieldArray []string) {
	ai := m.AddIndex{TableName: tableName}
	di := m.DropIndex{TableName: tableName}
	for key, value := range fieldArray {
		fieldArray[key] = strings.Trim(value, ", ")
		r, _ := regexp.Compile(FieldDataTypeRegexp)
		col := m.Columns{}
		if r.MatchString(fieldArray[key]) == true {
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
			ai = m.AddIndex{}
			di = m.DropIndex{}
		}
	}
	mUp.AddIndex = append(mUp.AddIndex, ai)
	mDown.DropIndex = append(mDown.DropIndex, di)
}

func fnChangeColumn(mUp *m.UpDown, mDown *m.UpDown) {
}

func fnCreateTable(mm *m.UpDown, tableName string, fieldArray []string) {
	ct := m.CreateTable{TableName: tableName}
	for key, value := range fieldArray {
		fieldArray[key] = strings.Trim(value, ", ")
		r, _ := regexp.Compile(FieldDataTypeRegexp)
		if r.MatchString(fieldArray[key]) == true {
			split := r.FindAllStringSubmatch(fieldArray[key], -1)
			col := m.Columns{
				FieldName: split[0][1],
				DataType:  split[0][2]}
			ct.Columns = append(ct.Columns, col)
		}
	}
	mm.CreateTable = append(mm.CreateTable, ct)
}

func fnDropColumn(mm *m.UpDown, tableName string, fieldArray []string) {
	dc := m.DropColumn{TableName: tableName}
	for key, value := range fieldArray {
		fieldArray[key] = strings.Trim(value, ", ")
		r, _ := regexp.Compile(FieldDataTypeRegexp)
		col := m.Columns{}
		if r.MatchString(fieldArray[key]) == true {
			split := r.FindAllStringSubmatch(fieldArray[key], -1)
			col.FieldName = split[0][1]
			col.DataType = split[0][2]
			dc.Columns = append(dc.Columns, col)
		} else if fieldArray[key] != "" {
			col.FieldName = fieldArray[key]
			dc.Columns = append(dc.Columns, col)
		}
	}
	mm.DropColumn = append(mm.DropColumn, dc)
}

func fnDropTable(mm *m.UpDown, tableName string) {
	dt := m.DropTable{TableName: tableName}
	mm.DropTable = append(mm.DropTable, dt)
}

func fnRenameTable(mUp *m.UpDown, mDown *m.UpDown, srcTableName string, destTablename string) {
	mUp.RenameTable = append(mUp.RenameTable, m.RenameTable{
		OldTableName: srcTableName,
		NewTableName: destTablename})
	mDown.RenameTable = append(mDown.RenameTable, m.RenameTable{
		OldTableName: destTablename,
		NewTableName: srcTableName})
}

// migrateUpDown - used to perform the Action Migration either Up or Down & chooses the DB too.
func migrateUpDown(updown string) {

	// Actual migration can't happen without having DB Name & DB User, rest can be default.
	if config.DbName == "" || config.DbUsername == "" {
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
			mm  m.Migration
			mig m.UpDown
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
			gdm_my.Init(config)
			if updown == "down" && mm.ID > gdm_my.GetLastMigrationNo() {
				continue
			}
			gdm_my.ProcessNow(mm, mig, updown, force)
		} else {
			gdm_pq.Init(config)
			if updown == "down" && mm.ID > gdm_my.GetLastMigrationNo() {
				continue
			}
			gdm_pq.ProcessNow(mm, mig, updown, force)
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
func checkConfigFileExists(c *m.Config) {
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
