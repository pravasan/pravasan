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
	currentVersion = "0.1"
	layout         = "20060102150405"

	// FieldDataTypeRegexp contains Regular Expression to split field name & field data type.
	FieldDataTypeRegexp = `^([A-Za-z_0-9$]{2,15}):([A-Za-z]{2,15})`
)

var (
	argArray              []string
	config                m.Config
	currentConfFileFormat = "json"
)

func main() {
	if strings.LastIndex(os.Args[0], "pravasan") < 1 || len(argArray) == 0 {
		fmt.Println("wrong usage, or no arguments specified")
		os.Exit(1)
	}
	switch argArray[0] {
	case "add", "a":
		generateMigration()
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
		panic("No or Wrong Actions provided.")
	}
	os.Exit(1)
}

func init() {
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
		migFileExtn     string
		migFilePrefix   string
		migOutputFormat string
		migTableName    string
	)

	flag.BoolVar(&flagPassword, "p", false, "database password")
	flag.BoolVar(&version, "version", false, "print Pravasan version")
	flag.StringVar(&configOutput, "confOutput", currentConfFileFormat, "config file format: json, xml")
	flag.StringVar(&dbHostname, "h", "localhost", "database hostname")
	flag.StringVar(&dbName, "d", "", "database name")
	flag.StringVar(&dbPort, "port", "5432", "database port")
	flag.StringVar(&dbType, "dbType", "mysql", "database type")
	flag.StringVar(&dbUsername, "u", "", "database username")
	flag.StringVar(&indexPrefix, "indexPrefix", "idx", "prefix for creating Indexes")
	flag.StringVar(&indexSuffix, "indexSuffix", "", "suffix for creating Indexes")
	flag.StringVar(&migFileExtn, "migFileExtn", "prvsn", "migration file extension")
	flag.StringVar(&migFilePrefix, "migFilePrefix", "", "prefix for migration file")
	flag.StringVar(&migOutputFormat, "migOutput", "json", "current supported format: json, xml")
	flag.StringVar(&migTableName, "migTableName", "schema_migrations", "migration table name")
	flag.Parse()

	if version {
		printCurrentVersion()
		if len(flag.Args()) == 0 {
			os.Exit(1)
		}
	}
	if flagPassword {
		fmt.Printf("Enter DB Password : ")
		pw := gopass.GetPasswd()
		config.DbPassword = string(pw)
	}

	if dbHostname != "" {
		config.DbHostname = dbHostname
	}
	if dbName != "" {
		config.DbName = dbName
	}
	if dbPort != "" {
		config.DbPort = dbPort
	}
	if dbType != "" {
		config.DbType = dbType
	}
	if dbUsername != "" {
		config.DbUsername = dbUsername
	}
	if indexPrefix != "" {
		config.IndexPrefix = indexPrefix
	}
	if indexSuffix != "" {
		config.IndexSuffix = indexSuffix
	}
	if migFileExtn != "" {
		config.MigrationFileExtension = migFileExtn
	}
	if migFilePrefix != "" {
		config.MigrationFilePrefix = migFilePrefix
	}
	if migOutputFormat != "" {
		config.MigrationOutputFormat = strings.ToLower(migOutputFormat)
	}
	if migTableName != "" {
		config.MigrationTableName = strings.ToLower(migTableName)
	}
	argArray = flag.Args()

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

func generateMigration() {
	t := time.Now()
	mm := m.Migration{}
	mm.ID = config.MigrationFilePrefix + t.Format(layout)
	switch argArray[1] {
	case "add_column", "ac":
		fnAddColumn(&mm.Up)
		fnDropColumn(&mm.Down)
	case "add_index", "ai":
		fnAddIndex(&mm.Up, &mm.Down)
	case "create_table", "ct":
		fnCreateTable(&mm.Up)
		fnDropTable(&mm.Down)
	case "drop_column", "dc":
		fnDropColumn(&mm.Up)
		fnAddColumn(&mm.Down)
	case "drop_index", "di":
		fnAddIndex(&mm.Down, &mm.Up)
	case "drop_table", "dt":
		fnDropTable(&mm.Up)
		fnCreateTable(&mm.Down)
	case "rename_table", "rt":
		fnRenameTable(&mm.Up, &mm.Down)
	case "sql", "s":
		fnSql(&mm)
	// case "change_column", "cc":
	// 	fnChangeColumn(&mm.Up, &mm.Down)
	// case "rename_column", "rc":
	// 	fnRenameColumn(&mm.Up, &mm.Down)
	default:
		panic("No or wrong Actions provided.")
	}

	writeToFile(mm.ID+"."+config.MigrationOutputFormat+"."+config.MigrationFileExtension, mm, config.MigrationOutputFormat)

	os.Exit(1)
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

func fnAddColumn(mm *m.UpDown) {
	ac := m.AddColumn{}
	ac.TableName = argArray[2]
	fieldArray := argArray[3:len(argArray)]
	for key, value := range fieldArray {
		fieldArray[key] = strings.Trim(value, ", ")
		r, _ := regexp.Compile(FieldDataTypeRegexp)
		if r.MatchString(fieldArray[key]) == true {
			split := r.FindAllStringSubmatch(fieldArray[key], -1)
			col := m.Columns{}
			col.FieldName = split[0][1]
			col.DataType = split[0][2]
			ac.Columns = append(ac.Columns, col)
		} else {
			ac = m.AddColumn{}
		}
	}
	mm.AddColumn = append(mm.AddColumn, ac)
}

func fnAddIndex(mUp *m.UpDown, mDown *m.UpDown) {
	ai := m.AddIndex{}
	di := m.DropIndex{}
	ai.TableName = argArray[2]
	di.TableName = argArray[2]
	fieldArray := argArray[3:len(argArray)]
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

func fnCreateTable(mm *m.UpDown) {
	ct := m.CreateTable{}
	ct.TableName = argArray[2]
	fieldArray := argArray[3:len(argArray)]
	for key, value := range fieldArray {
		fieldArray[key] = strings.Trim(value, ", ")
		r, _ := regexp.Compile(FieldDataTypeRegexp)
		if r.MatchString(fieldArray[key]) == true {
			split := r.FindAllStringSubmatch(fieldArray[key], -1)
			col := m.Columns{}
			col.FieldName = split[0][1]
			col.DataType = split[0][2]
			ct.Columns = append(ct.Columns, col)
		}
	}
	mm.CreateTable = append(mm.CreateTable, ct)
}

func fnDropColumn(mm *m.UpDown) {
	dc := m.DropColumn{}
	dc.TableName = argArray[2]
	fieldArray := argArray[3:len(argArray)]
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
		} else {
			dc = m.DropColumn{}
		}
	}
	mm.DropColumn = append(mm.DropColumn, dc)
}

func fnDropTable(mm *m.UpDown) {
	dt := m.DropTable{}
	dt.TableName = argArray[2]
	mm.DropTable = append(mm.DropTable, dt)
}

func fnRenameTable(mUp *m.UpDown, mDown *m.UpDown) {
	rt := m.RenameTable{}
	rt.OldTableName = argArray[2]
	rt.NewTableName = argArray[3]
	mUp.RenameTable = append(mUp.RenameTable, rt)
	rt.OldTableName = argArray[3]
	rt.NewTableName = argArray[2]
	mDown.RenameTable = append(mDown.RenameTable, rt)
}

// migrateUpDown - used to perform the Action Migration either Up or Down & chooses the DB too.
func migrateUpDown(updown string) {

	// Actual migration can't happen without having DB Name & DB User, rest can be default.
	if config.DbName == "" || config.DbUsername == "" {
		fmt.Println("Either Database Name or Username is not mentioned, or both are missed to mention.")
		return
	}

	var (
		err          error
		files        []string
		processCount int
		reverseCount = 1
		force        = false
	)

	files = migrationFiles(updown)
	if 0 == len(files) {
		fmt.Println("No files in the directory")
		return
	}

	if len(argArray) > 1 && argArray[1] != "" {
		if updown == "down" {
			reverseCount, err = strconv.Atoi(strings.Replace(argArray[1], "-", "", -1))
			if err != nil {
				log.Println("Wrong count to be reversed, only integer values accepted")
				log.Fatal(err)
			}
		} else if updown == "up" {
			files, err = checkMigrationFilesExists(argArray[1:len(argArray)])
			if err != nil {
				fmt.Println("one of the version number of the file mentioned is wrong.")
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
		fmt.Println("Processing ... " + filename)
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
			err = errors.New("coudn't able to read file: " + indvlFile)
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
	files, _ := ioutil.ReadDir("./")
	var onlyMigFiles []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), "."+config.MigrationOutputFormat+"."+config.MigrationFileExtension) && strings.HasPrefix(f.Name(), config.MigrationFilePrefix) {
			onlyMigFiles = append(onlyMigFiles, f.Name())
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
	fmt.Println("Config file created.")
}

func printCurrentVersion() {
	fmt.Println("pravasan version " + currentVersion)
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
