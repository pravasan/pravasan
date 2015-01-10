package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
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
	// FieldDataTypeRegexp contains Regular Expression to split field name & field data type.
	FieldDataTypeRegexp = `^([A-Za-z_0-9$]{2,15}):([A-Za-z]{2,15})`
	currentVersion      = "0.1"
	layout              = "20060102150405"
)

var (
	config   m.Config
	argArray []string
)

func main() {
	if strings.LastIndex(os.Args[0], "pravasan") < 1 || len(argArray) == 0 {
		fmt.Println("wrong usage, or no arguments specified")
		os.Exit(1)
	}
	switch argArray[0] {
	case "add", "a":
		generateMigration()
	case "up", "u":
		migrateUpDown("up")
	case "down", "d":
		migrateUpDown("down")
	case "create", "c":
		if len(argArray) > 1 && argArray[1] != "" && argArray[1] == "conf" {
			createConfigurationFile()
		} else {
			createMigration()
		}
	default:
		panic("No or Wrong Actions provided.")
	}
	os.Exit(1)
}

func createConfigurationFile() {
	writeToFile("pravasan.conf."+config.MigrationOutput, config, config.MigrationOutput)
	fmt.Println("Config file created / updated.")
}

func init() {
	var currentConfFileFormat = "json"
	// fmt.Println("pravasan init() it runs before other functions")
	if _, err := os.Stat("./pravasan.conf.json"); err == nil {
		currentConfFileFormat = "json"
		bs, err := ioutil.ReadFile("pravasan.conf.json")
		if err != nil {
			fmt.Println(err)
			return
		}
		docScript := []byte(bs)
		json.Unmarshal(docScript, &config)
	} else if _, err := os.Stat("./pravasan.conf.xml"); err == nil {
		currentConfFileFormat = "xml"
		bs, err := ioutil.ReadFile("pravasan.conf.xml")
		if err != nil {
			fmt.Println(err)
			return
		}
		docScript := []byte(bs)
		xml.Unmarshal(docScript, &config)
	}

	var (
		dbname             string
		dbType             string
		extn               string
		flagPassword       = false
		host               string
		migrationTableName string
		output             string
		port               string
		prefix             string
		un                 string
		version            = false
	)

	flag.BoolVar(&flagPassword, "p", false, "specify the option asking for database password")
	flag.BoolVar(&version, "version", false, "print Pravasan version")
	flag.StringVar(&dbname, "d", "", "specify the database name")
	flag.StringVar(&dbType, "dbType", "mysql", "specify the database type")
	flag.StringVar(&extn, "extn", "prvsn", "specify the migration file extension")
	flag.StringVar(&host, "h", "localhost", "specify the database hostname")
	flag.StringVar(&migrationTableName, "migration_table_name", "schema_migrations", "supported format are json, xml")
	flag.StringVar(&output, "output", currentConfFileFormat, "supported format are json, xml")
	flag.StringVar(&port, "port", "5432", "specify the database port")
	flag.StringVar(&prefix, "prefix", "", "specify the text to be prefix with the migration file")
	flag.StringVar(&un, "u", "", "specify the database username")
	flag.Parse()

	if version {
		printCurrentVersion()
		if len(flag.Args()) == 0 {
			os.Exit(1)
		}
	}

	if dbType != "" {
		config.DbType = dbType
	}
	if un != "" {
		config.DbUsername = un
	}
	if flagPassword {
		fmt.Printf("Enter DB Password : ")
		pw := gopass.GetPasswd()
		config.DbPassword = string(pw)
	}
	if dbname != "" {
		config.DbName = dbname
	}
	if host != "" {
		config.DbHostname = host
	}
	if port != "" {
		config.DbPortnumber = port
	}
	if prefix != "" {
		config.FilePrefix = prefix
	}
	if extn != "" {
		config.FileExtension = extn
	}
	if output != "" {
		config.MigrationOutput = strings.ToLower(output)
	}
	if migrationTableName != "" {
		config.MigrationTableName = strings.ToLower(migrationTableName)
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
	mm.ID = config.FilePrefix + t.Format(layout)
	switch argArray[1] {
	case "create_table", "ct":
		fnCreateTable(&mm.Up)
		fnDropTable(&mm.Down)
	case "drop_table", "dt":
		fnDropTable(&mm.Up)
		fnCreateTable(&mm.Down)
	case "rename_table", "rt":
		fnRenameTable(&mm.Up, &mm.Down)
	case "add_column", "ac":
		fnAddColumn(&mm.Up)
		fnDropColumn(&mm.Down)
	case "drop_column", "dc":
		fnDropColumn(&mm.Up)
		fnAddColumn(&mm.Down)
	case "sql", "s":
		fnSql(&mm)
	// case "change_column", "cc":
	// 	fnChangeColumn(&mm.Up, &mm.Down)
	// case "rename_column", "rc":
	// 	fnRenameColumn(&mm.Up, &mm.Down)
	// case "add_index", "ai":
	// 	fnAddIndex(&mm.Up, &mm.Down)
	default:
		panic("No or wrong Actions provided.")
	}

	writeToFile(mm.ID+"."+config.MigrationOutput+"."+config.FileExtension, mm, config.MigrationOutput)

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

	fmt.Println("Enter SQL statements for Down section of migration")
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

func fnAddIndex(mUp *m.UpDown, mDown *m.UpDown) {
	// ct := m.CreateTable{}
	// ct.TableName = argArray[2]
	// fieldArray := argArray[3:len(argArray)]
	// for key, value := range fieldArray {
	// 	fieldArray[key] = strings.Trim(value, ", ")
	// 	r, _ := regexp.Compile(FieldDataTypeRegexp)
	// 	if r.MatchString(fieldArray[key]) == true {
	// 		split := r.FindAllStringSubmatch(fieldArray[key], -1)
	// 		col := m.Columns{}
	// 		col.FieldName = split[0][1]
	// 		col.DataType = split[0][2]
	// 		ct.Columns = append(ct.Columns, col)
	// 	}
	// }
	// mm.CreateTable = append(mm.CreateTable, ct)

	// ai := m.AddIndex{}
	// ai.TableName = argArray[2]
	// fieldArray = argArray[3:len(argArray)]
	// for key, value := range fieldArray {
	// 	fieldArray[key] = strings.Trim(value, ", ")
	// 	r, _ := regexp.Compile(FieldDataTypeRegexp)
	// 	if r.MatchString(fieldArray[key]) == true {
	// 		split := r.FindAllStringSubmatch(fieldArray[key], -1)
	// 		col := m.Columns{}
	// 		col.FieldName = split[0][1]
	// 		col.DataType = split[0][2]
	// 		ai.Columns = append(ai.Columns, col)
	// 	}
	// }
	// mUp.AddIndex = append(mUp.AddIndex, ai)
}

func fnChangeColumn(mUp *m.UpDown, mDown *m.UpDown) {
}

func fnRenameColumn(mUp *m.UpDown, mDown *m.UpDown) {
}

func fnRenameTable(mUp *m.UpDown, mDown *m.UpDown) {

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

// migrateUpDown - used to perform the Action Migration either Up or Down & chooses the DB too.
func migrateUpDown(updown string) {
	if config.DbName == "" || config.DbUsername == "" {
		fmt.Println("Either Database Name or Username is not mentioned, or both are missed to mention.")
	}

	files := migrationFiles(updown)
	if 0 == len(files) {
		fmt.Println("No files in the directory")
		return
	}

	var processCount int

	// setting reverseCount
	var reverseCount = 1
	var err error
	if len(argArray) > 1 && argArray[1] != "" {
		reverseCount, err = strconv.Atoi(strings.Replace(argArray[1], "-", "", -1))
		if err != nil {
			log.Println("Wrong count to be reversed, only integer values accepted")
			log.Fatal(err)
		}
	}

	for _, filename := range files {

		// During migration down if count is reached then exit
		if updown == "down" && reverseCount == processCount {
			break
		}

		// Read the content of the file & store into the mm structure
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

		if config.MigrationOutput == "json" {
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
			gdm_my.ProcessNow(mm, mig, updown)
		} else {
			gdm_pq.Init(config)
			if updown == "down" && mm.ID > gdm_my.GetLastMigrationNo() {
				continue
			}
			gdm_pq.ProcessNow(mm, mig, updown)
		}
		processCount++
	}
}

func migrationFiles(updown string) []string {
	files, _ := ioutil.ReadDir("./")
	var onlyMigFiles []string
	for _, f := range files {
		if !f.IsDir() && strings.Contains(f.Name(), "."+config.MigrationOutput+"."+config.FileExtension) {
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

func printCurrentVersion() {
	fmt.Println("pravasan version " + currentVersion)
}
