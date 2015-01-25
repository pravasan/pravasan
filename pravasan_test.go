package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestPrintCurrentVersion(t *testing.T) {
	pcv := printCurrentVersion()
	if pcv != infoText+"pravasan version 0.3"+resetText {
		t.Error("Failed to identify the correct version")
	}
}

func Test_AddColumn_Migration(t *testing.T) {
	argsss := []string{"add", "ac", "test123", "new_column:int"}
	fileName, mm := generateMigration(argsss)
	expectedString := `{"id":"` + strings.TrimSuffix(fileName, "..") + `","up":{"addColumn":[{"tableName":"test123","columns":[{"fieldname":"new_column","datatype":"int"}]}]},"down":{"dropColumn":[{"tableName":"test123","columns":[{"fieldname":"new_column","datatype":"int"}]}]}}`
	content1, _ := json.Marshal(mm)
	checkError(t, expectedString, string(content1))
}

func Test_DropColumn_Migration(t *testing.T) {
	argsss := []string{"add", "dc", "test123", "new_column"}
	fileName, mm := generateMigration(argsss)
	expectedString := `{"id":"` + strings.TrimSuffix(fileName, "..") + `","up":{"dropColumn":[{"tableName":"test123","columns":[{"fieldname":"new_column"}]}]},"down":{"addColumn":[{}]}}`
	content1, _ := json.Marshal(mm)
	checkError(t, expectedString, string(content1))
}

func Test_CreateTable_Migration(t *testing.T) {
	argsss := []string{"add", "ct", "test123", "first:int", "second:string"}
	fileName, mm := generateMigration(argsss)
	expectedString := `{"id":"` + strings.TrimSuffix(fileName, "..") + `","up":{"createTable":[{"tableName":"test123","columns":[{"fieldname":"first","datatype":"int"},{"fieldname":"second","datatype":"string"}]}]},"down":{"dropTable":[{"tableName":"test123"}]}}`
	content1, _ := json.Marshal(mm)
	checkError(t, expectedString, string(content1))
}

func Test_DropTable_Migration(t *testing.T) {
	argsss := []string{"add", "dt", "test123"}
	fileName, mm := generateMigration(argsss)
	expectedString := `{"id":"` + strings.TrimSuffix(fileName, "..") + `","up":{"dropTable":[{"tableName":"test123"}]},"down":{"createTable":[{"tableName":"test123"}]}}`
	content1, _ := json.Marshal(mm)
	checkError(t, expectedString, string(content1))
}

func Test_RenameTable_Migration(t *testing.T) {
	argsss := []string{"add", "rt", "old_test123", "new_test123"}
	fileName, mm := generateMigration(argsss)
	expectedString := `{"id":"` + strings.TrimSuffix(fileName, "..") + `","up":{"renameTable":[{"oldTableName":"old_test123","newTableName":"new_test123"}]},"down":{"renameTable":[{"oldTableName":"new_test123","newTableName":"old_test123"}]}}`
	content1, _ := json.Marshal(mm)
	checkError(t, expectedString, string(content1))
}

func Test_AddIndex_Migration(t *testing.T) {
	argsss := []string{"add", "ai", "test123", "first_col", "second_col", "third_col"}
	fileName, mm := generateMigration(argsss)
	expectedString := `{"id":"` + strings.TrimSuffix(fileName, "..") + `","up":{"addIndex":[{"tableName":"test123","columns":[{"fieldname":"first_col"},{"fieldname":"second_col"},{"fieldname":"third_col"}]}]},"down":{"dropIndex":[{"tableName":"test123","columns":[{"fieldname":"first_col"},{"fieldname":"second_col"},{"fieldname":"third_col"}]}]}}`
	content1, _ := json.Marshal(mm)
	checkError(t, expectedString, string(content1))
}

func Test_DropIndex_Migration(t *testing.T) {
	argsss := []string{"add", "di", "test123", "first_col", "second_col", "third_col"}
	fileName, mm := generateMigration(argsss)
	expectedString := `{"id":"` + strings.TrimSuffix(fileName, "..") + `","up":{"dropIndex":[{"tableName":"test123","columns":[{"fieldname":"first_col"},{"fieldname":"second_col"},{"fieldname":"third_col"}]}]},"down":{"addIndex":[{"tableName":"test123","columns":[{"fieldname":"first_col"},{"fieldname":"second_col"},{"fieldname":"third_col"}]}]}}`
	content1, _ := json.Marshal(mm)
	checkError(t, expectedString, string(content1))
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
