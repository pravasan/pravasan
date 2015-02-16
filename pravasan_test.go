package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestPrintCurrentVersion(t *testing.T) {
	pcv := printCurrentVersion()
	if pcv != infoText+"pravasan version 0.4"+resetText {
		t.Error("Failed to identify the correct version")
	}
}

func Test_CurrentVersion_Using_Flag(t *testing.T) {
	// // #TODO(kishorevaishnav): Don't know how to handle this
	// args := os.Args
	// defer func() { os.Args = args }()
	// os.Args = []string{"", "-version"}
	// initializeDefaults()
}

func Test_AddColumn_Migration(t *testing.T) {
	argsss := []string{"add", "ac", "test123", "new_column:int"}
	fileName, mm, _ := generateMigration(argsss)
	expectedString := `{"id":"` + getID(fileName) + `","up":{"addColumn":[{"tableName":"test123","columns":[{"fieldname":"new_column","datatype":"int"}]}]},"down":{"dropColumn":[{"tableName":"test123","columns":[{"fieldname":"new_column","datatype":"int"}]}]}}`
	content1, _ := json.Marshal(mm)
	checkError(t, expectedString, string(content1))
}

func Test_DropColumn_Migration(t *testing.T) {
	argsss := []string{"add", "dc", "test123", "new_column"}
	fileName, mm, _ := generateMigration(argsss)
	expectedString := `{"id":"` + getID(fileName) + `","up":{"dropColumn":[{"tableName":"test123","columns":[{"fieldname":"new_column"}]}]},"down":{"addColumn":[{"tableName":"test123","columns":[{"fieldname":"new_column"}]}]}}`
	content1, _ := json.Marshal(mm)
	checkError(t, expectedString, string(content1))
}

func Test_CreateTable_Migration(t *testing.T) {
	argsss := []string{"add", "ct", "test123", "first:int", "second:string"}
	fileName, mm, _ := generateMigration(argsss)
	expectedString := `{"id":"` + getID(fileName) + `","up":{"createTable":[{"tableName":"test123","columns":[{"fieldname":"first","datatype":"int"},{"fieldname":"second","datatype":"string"}]}]},"down":{"dropTable":[{"tableName":"test123"}]}}`
	content1, _ := json.Marshal(mm)
	checkError(t, expectedString, string(content1))
}

func Test_CreateTableWithSize_Migration(t *testing.T) {
	argsss := []string{"add", "ct", "test123", "first:int", "second:varchar(255)"}
	fileName, mm, _ := generateMigration(argsss)
	expectedString := `{"id":"` + getID(fileName) + `","up":{"createTable":[{"tableName":"test123","columns":[{"fieldname":"first","datatype":"int"},{"fieldname":"second","datatype":"varchar(255)"}]}]},"down":{"dropTable":[{"tableName":"test123"}]}}`
	content1, _ := json.Marshal(mm)
	checkError(t, expectedString, string(content1))
}

func Test_CreateTable_Migration_AutoAddColumns(t *testing.T) {
	// #TODO(kishorevaishnav): There is some problem with the below code.
	f1 := flag.NewFlagSet("f1", flag.ContinueOnError)
	var autoAddColumns string
	f1.StringVar(&autoAddColumns, "autoAddColumns", "id:int created_at:datetime modified_at:datetime", "-------")
	err := f1.Parse([]string{"autoAddColumns"})
	if err != nil {
		fmt.Println(autoAddColumns)
		fmt.Println(err)
	}
	// initializeDefaults()
	argsss := []string{"add", "ct", "test123", "first:int", "second:string"}
	fileName, mm, _ := generateMigration(argsss)
	expectedString := `{"id":"` + getID(fileName) + `","up":{"createTable":[{"tableName":"test123","columns":[{"fieldname":"first","datatype":"int"},{"fieldname":"second","datatype":"string"}]}]},"down":{"dropTable":[{"tableName":"test123"}]}}`
	content1, _ := json.Marshal(mm)
	checkError(t, expectedString, string(content1))
}

func Test_DropTable_Migration(t *testing.T) {
	argsss := []string{"add", "dt", "test123"}
	fileName, mm, _ := generateMigration(argsss)
	expectedString := `{"id":"` + getID(fileName) + `","up":{"dropTable":[{"tableName":"test123"}]},"down":{"createTable":[{"tableName":"test123"}]}}`
	content1, _ := json.Marshal(mm)
	checkError(t, expectedString, string(content1))
}

func Test_RenameTable_Migration(t *testing.T) {
	argsss := []string{"add", "rt", "old_test123", "new_test123"}
	fileName, mm, _ := generateMigration(argsss)
	expectedString := `{"id":"` + getID(fileName) + `","up":{"renameTable":[{"oldTableName":"old_test123","newTableName":"new_test123"}]},"down":{"renameTable":[{"oldTableName":"new_test123","newTableName":"old_test123"}]}}`
	content1, _ := json.Marshal(mm)
	checkError(t, expectedString, string(content1))
}

func Test_AddIndex_Migration(t *testing.T) {
	argsss := []string{"add", "ai", "test123", "first_col", "second_col", "third_col"}
	fileName, mm, _ := generateMigration(argsss)
	expectedString := `{"id":"` + getID(fileName) + `","up":{"addIndex":[{"tableName":"test123","columns":[{"fieldname":"first_col"},{"fieldname":"second_col"},{"fieldname":"third_col"}]}]},"down":{"dropIndex":[{"tableName":"test123","columns":[{"fieldname":"first_col"},{"fieldname":"second_col"},{"fieldname":"third_col"}]}]}}`
	content1, _ := json.Marshal(mm)
	checkError(t, expectedString, string(content1))
}

func Test_DropIndex_Migration(t *testing.T) {
	argsss := []string{"add", "di", "test123", "first_col", "second_col", "third_col"}
	fileName, mm, _ := generateMigration(argsss)
	expectedString := `{"id":"` + getID(fileName) + `","up":{"dropIndex":[{"tableName":"test123","columns":[{"fieldname":"first_col"},{"fieldname":"second_col"},{"fieldname":"third_col"}]}]},"down":{"addIndex":[{"tableName":"test123","columns":[{"fieldname":"first_col"},{"fieldname":"second_col"},{"fieldname":"third_col"}]}]}}`
	content1, _ := json.Marshal(mm)
	checkError(t, expectedString, string(content1))
}

func Test_WithoutConfig_MigrationDirectoryExists(t *testing.T) {
	args := os.Args
	defer func() { os.Args = args }()
	os.Args = []string{""}
	initializeDefaults()
	_, err := migrationDirectoryExists()
	if err != nil {
		t.Error("failed....")
		t.Log(err)
	}
}

func Test_WithConfig_MigrationDirectoryExists(t *testing.T) {
	// #TODO(kishorevaishnav): pending to test this function.
	args := os.Args
	defer func() { os.Args = args }()
	// os.Args = []string{"./pravasan", "add", "ct", "test123", "first:int"}
	os.Args = []string{"", "-migDir=migrate"}

	flagSet := flag.NewFlagSet("example", flag.ContinueOnError)
	migDira := flagSet.String("migDir", "migrate", "filename to view")
	err := flagSet.Parse([]string{"migDir"})
	if err != nil {
		fmt.Println(*migDira)
		fmt.Println(err)
	}
	initializeDefaults()
	_, err = migrationDirectoryExists()
	if err == nil {
		t.Error("expected to fail....")
		t.Log(err)
	}
}

func TestMain(t *testing.T) {
	args := os.Args
	defer func() { os.Args = args }()
	// os.Args = []string{"./pravasan", "add", "ct", "test123", "first:int"}
	// fmt.Println(os.Args)
	// main()
}

func checkError(t *testing.T, expString string, genString string) {
	if expString != genString {
		reportError(t, expString, genString)
	}
}

func reportError(t *testing.T, expOutput string, genOutput string) {
	fmt.Println("Expected Output")
	fmt.Println(expOutput)
	fmt.Println("\nGenerated Output")
	fmt.Println(genOutput)
	t.Error("DOESN'T MATCH")
}

func getID(fileName string) string {
	return strings.TrimSuffix(strings.TrimPrefix(strings.TrimSuffix(fileName, ".json.prvsn"), "./"), "..")
}
