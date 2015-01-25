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
	if pcv == infoText+"pravasan version 0.2"+resetText {
		t.Log("Success")
	} else {
		t.Error("Failed to identify the correct version")
	}
}
func TestWriteToFile(t *testing.T) {
	t.Log("Success")
}

func TestAnotherScenario(t *testing.T) {
	t.Log("Success")
}

func Test_CreateTable_Migration(t *testing.T) {
	argsss := []string{"add", "ct", "test123", "first:int", "second:string"}
	fileName, mm := generateMigration(argsss)
	expectedString := `{"id":"` + strings.TrimSuffix(fileName, "..") + `","up":{"createTable":[{"tableName":"test123","columns":[{"fieldname":"first","datatype":"int"},{"fieldname":"second","datatype":"string"}]}]},"down":{"dropTable":[{"tableName":"test123"}]}}`
	content1, _ := json.Marshal(mm)
	if expectedString != string(content1) {
		fmt.Println("Expected Output")
		fmt.Println(expectedString)
		fmt.Println("\nGenerated Output")
		fmt.Println(string(content1))
		t.Error("Doesn't match")
	}
	t.Log("Success")
}

// func Test_DropTable_Migration(t *testing.T) {
// 	argsss := []string{"add", "ct", "test123", "first:int", "second:string"}
// 	fileName, mm := generateMigration(argsss)
// 	expectedString := `{"id":"` + strings.TrimSuffix(fileName, "..") + `","up":{"createTable":[{"tableName":"test123","columns":[{"fieldname":"first","datatype":"int"},{"fieldname":"second","datatype":"string"}]}]},"down":{"dropTable":[{"tableName":"test123"}]}}`
// 	content1, _ := json.Marshal(mm)
// 	if expectedString != string(content1) {
// 		fmt.Println("Expected Output")
// 		fmt.Println(expectedString)
// 		fmt.Println("\nGenerated Output")
// 		fmt.Println(string(content1))
// 		t.Error("Doesn't match")
// 	}
// 	t.Log("Success")
// }

func TestMain(t *testing.T) {
	args := os.Args
	defer func() { os.Args = args }()
	// os.Args = []string{"./pravasan", "add", "ct", "test123", "first:int"}
	// fmt.Println(os.Args)
	// main()
	t.Log("Success")
}
