package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
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
	FIELD_DATATYPE_REGEXP = `^([A-Za-z_0-9$]{2,15}):([A-Za-z]{2,15})`
	current_version       = "0.1"
	layout                = "20060102150405"
)

var (
	Config m.Config
	ArgArr []string
)

func main() {
	if strings.LastIndex(os.Args[0], "pravasan") < 1 || len(ArgArr) == 0 {
		fmt.Println("wrong usage, or no arguments specified")
		os.Exit(1)
	}
	switch ArgArr[0] {
	case "add", "a":
		generateMigration()
	case "up", "u":
		migrateUpDown("up")
	case "down", "d":
		migrateUpDown("down")
	case "create", "c":
		if len(ArgArr) > 1 && ArgArr[1] != "" && ArgArr[1] == "conf" {
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
	writeToFile("pravasan.conf."+Config.Migration_output, Config, Config.Migration_output)
	fmt.Println("Config file created / updated.")
}

func init() {
	var current_conf_file_format string = "json"
	// fmt.Println("pravasan init() it runs before other functions")
	if _, err := os.Stat("./pravasan.conf.json"); err == nil {
		current_conf_file_format = "json"
		bs, err := ioutil.ReadFile("pravasan.conf.json")
		if err != nil {
			fmt.Println(err)
			return
		}
		docScript := []byte(bs)
		json.Unmarshal(docScript, &Config)
	} else if _, err := os.Stat("./pravasan.conf.xml"); err == nil {
		current_conf_file_format = "xml"
		bs, err := ioutil.ReadFile("pravasan.conf.xml")
		if err != nil {
			fmt.Println(err)
			return
		}
		docScript := []byte(bs)
		xml.Unmarshal(docScript, &Config)
	}

	var (
		db_type       string
		un            string
		dbname        string
		host          string
		port          string
		prefix        string
		extn          string
		output        string
		version       bool = false
		flag_password bool = false
	)

	flag.StringVar(&db_type, "db_type", "mysql", "specify the database type")
	flag.StringVar(&un, "u", "", "specify the database username")
	flag.BoolVar(&flag_password, "p", false, "specify the option asking for database password")
	flag.StringVar(&dbname, "d", "", "specify the database name")
	flag.StringVar(&host, "h", "localhost", "specify the database hostname")
	flag.StringVar(&port, "port", "5432", "specify the database port")
	flag.StringVar(&prefix, "prefix", "", "specify the text to be prefix with the migration file")
	flag.StringVar(&extn, "extn", "prvsn", "specify the migration file extension")
	flag.BoolVar(&version, "version", false, "print Pravasan version")
	flag.StringVar(&output, "output", current_conf_file_format, "supported format are json, xml")
	flag.Parse()

	if version {
		printCurrentVersion()
		if len(flag.Args()) == 0 {
			os.Exit(1)
		}
	}

	if db_type != "" {
		Config.Db_type = db_type
	}
	if un != "" {
		Config.Db_username = un
	}
	if flag_password {
		fmt.Printf("Enter DB Password : ")
		pw := gopass.GetPasswd()
		Config.Db_password = string(pw)
	}
	if dbname != "" {
		Config.Db_name = dbname
	}
	if host != "" {
		Config.Db_hostname = host
	}
	if port != "" {
		Config.Db_portnumber = port
	}
	if prefix != "" {
		Config.File_prefix = prefix
	}
	if extn != "" {
		Config.File_extension = extn
	}
	if output != "" {
		Config.Migration_output = strings.ToLower(output)
	}
	ArgArr = flag.Args()

}

func createMigration() {
	if Config.Db_type == "mysql" {
		gdm_my.Init(Config)
		gdm_my.CreateMigrationTable()
	} else {
		gdm_pq.Init(Config)
		gdm_pq.CreateMigrationTable()
	}
}

func generateMigration() {
	t := time.Now()
	mm := m.Migration{}
	mm.Id = Config.File_prefix + t.Format(layout)
	switch ArgArr[1] {
	case "create_table", "ct":
		fn_create_table(&mm.Up)
		fn_drop_table(&mm.Down)
	case "drop_table", "dt":
		fn_drop_table(&mm.Up)
		fn_create_table(&mm.Down)
	case "rename_table", "rt":
		fn_rename_table(&mm.Up, &mm.Down)
	case "add_column", "ac":
		fn_add_column(&mm.Up)
		fn_drop_column(&mm.Down)
	case "drop_column", "dc":
		fn_drop_column(&mm.Up)
		fn_add_column(&mm.Down)
	// case "change_column", "cc":
	// 	fn_change_column(&mm.Up, &mm.Down)
	// case "rename_column", "rc":
	// 	fn_rename_column(&mm.Up, &mm.Down)
	// case "add_index", "ai":
	// 	fn_add_index(&mm.Up, &mm.Down)
	default:
		panic("No or wrong Actions provided.")
	}

	writeToFile(mm.Id+"."+Config.Migration_output+"."+Config.File_extension, mm, Config.Migration_output)

	os.Exit(1)
}

func fn_add_index(mm_up *m.UpDown, mm_down *m.UpDown) {
	// ct := m.CreateTable{}
	// ct.Table_Name = ArgArr[2]
	// fieldArray := ArgArr[3:len(ArgArr)]
	// for key, value := range fieldArray {
	// 	fieldArray[key] = strings.Trim(value, ", ")
	// 	r, _ := regexp.Compile(FIELD_DATATYPE_REGEXP)
	// 	if r.MatchString(fieldArray[key]) == true {
	// 		split := r.FindAllStringSubmatch(fieldArray[key], -1)
	// 		col := m.Columns{}
	// 		col.FieldName = split[0][1]
	// 		col.DataType = split[0][2]
	// 		ct.Columns = append(ct.Columns, col)
	// 	}
	// }
	// mm.Create_Table = append(mm.Create_Table, ct)

	// ai := m.AddIndex{}
	// ai.Table_Name = ArgArr[2]
	// fieldArray = ArgArr[3:len(ArgArr)]
	// for key, value := range fieldArray {
	// 	fieldArray[key] = strings.Trim(value, ", ")
	// 	r, _ := regexp.Compile(FIELD_DATATYPE_REGEXP)
	// 	if r.MatchString(fieldArray[key]) == true {
	// 		split := r.FindAllStringSubmatch(fieldArray[key], -1)
	// 		col := m.Columns{}
	// 		col.FieldName = split[0][1]
	// 		col.DataType = split[0][2]
	// 		ai.Columns = append(ai.Columns, col)
	// 	}
	// }
	// mm_up.Add_Index = append(mm_up.Add_Index, ai)
}

func fn_change_column(mm_up *m.UpDown, mm_down *m.UpDown) {
}

func fn_rename_column(mm_up *m.UpDown, mm_down *m.UpDown) {
}

func fn_rename_table(mm_up *m.UpDown, mm_down *m.UpDown) {

}

func fn_add_column(mm *m.UpDown) {
	ac := m.AddColumn{}
	ac.Table_Name = ArgArr[2]
	fieldArray := ArgArr[3:len(ArgArr)]
	for key, value := range fieldArray {
		fieldArray[key] = strings.Trim(value, ", ")
		r, _ := regexp.Compile(FIELD_DATATYPE_REGEXP)
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
	mm.Add_Column = append(mm.Add_Column, ac)
}

func fn_drop_column(mm *m.UpDown) {
	dc := m.DropColumn{}
	dc.Table_Name = ArgArr[2]
	fieldArray := ArgArr[3:len(ArgArr)]
	for key, value := range fieldArray {
		fieldArray[key] = strings.Trim(value, ", ")
		r, _ := regexp.Compile(FIELD_DATATYPE_REGEXP)
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
	mm.Drop_Column = append(mm.Drop_Column, dc)
}

func fn_drop_table(mm *m.UpDown) {
	dt := m.DropTable{}
	dt.Table_Name = ArgArr[2]
	mm.Drop_Table = append(mm.Drop_Table, dt)
}

func fn_create_table(mm *m.UpDown) {
	ct := m.CreateTable{}
	ct.Table_Name = ArgArr[2]
	fieldArray := ArgArr[3:len(ArgArr)]
	for key, value := range fieldArray {
		fieldArray[key] = strings.Trim(value, ", ")
		r, _ := regexp.Compile(FIELD_DATATYPE_REGEXP)
		if r.MatchString(fieldArray[key]) == true {
			split := r.FindAllStringSubmatch(fieldArray[key], -1)
			col := m.Columns{}
			col.FieldName = split[0][1]
			col.DataType = split[0][2]
			ct.Columns = append(ct.Columns, col)
		}
	}
	mm.Create_Table = append(mm.Create_Table, ct)
}

// migrateUpDown - used to perform the Action Migration either Up or Down & chooses the DB too.
func migrateUpDown(updown string) {
	if Config.Db_name == "" || Config.Db_username == "" {
		fmt.Println("Either Database Name or Username is not mentioned, or both are missed to mention.")
	}

	files := migrationFiles(updown)
	if 0 == len(files) {
		fmt.Println("No files in the directory")
		return
	}

	var processCount int = 0

	// setting reverseCount
	var reverseCount int = 1
	var err error
	if len(ArgArr) > 1 && ArgArr[1] != "" {
		reverseCount, err = strconv.Atoi(strings.Replace(ArgArr[1], "-", "", -1))
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

		if Config.Migration_output == "json" {
			json.Unmarshal(docScript, &mm)
		} else {
			xml.Unmarshal(docScript, &mm)
		}

		if updown == "up" {
			mig = mm.Up
		} else {
			mig = mm.Down
		}

		if Config.Db_type == "mysql" {
			gdm_my.Init(Config)
			if updown == "down" && mm.Id > gdm_my.GetLastMigrationNo() {
				continue
			}
			gdm_my.ProcessNow(mm, mig, updown)
		} else {
			gdm_pq.Init(Config)
			if updown == "down" && mm.Id > gdm_my.GetLastMigrationNo() {
				continue
			}
			gdm_pq.ProcessNow(mm, mig, updown)
		}
		processCount++
	}
}

func migrationFiles(updown string) []string {
	files, _ := ioutil.ReadDir("./*.json.prvsn")
	var onlyMigFiles []string
	for _, f := range files {
		if !f.IsDir() && strings.Contains(f.Name(), "."+Config.Migration_output+"."+Config.File_extension) {
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
	fmt.Println("pravasan version " + current_version)
}
