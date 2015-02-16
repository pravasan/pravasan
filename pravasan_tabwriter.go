package main

import (
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// TableOfDataTypes - prints  mentioned database or all databases datatypes in table format
func TableOfDataTypes(dbname string) {
	abc := [][]string{}
	for db, list := range ListSuppDataTypes {
		if dbname == "" || strings.ToUpper(dbname) == strings.ToUpper(db) {
			for k, v := range list {
				if k == v {
					v = ""
				}
				abc = append(abc, []string{k, db, v})
			}
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"DataType", "Supported Database(s)", "Alias"})

	for _, v := range abc {
		table.Append(v)
	}
	table.Render() // Send output
	return
}
