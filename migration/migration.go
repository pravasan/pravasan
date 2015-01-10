package migration

// Config used for general configuration settings for Migration application
type Config struct {
	DbType             string `xml:"DbType,omitempty" json:"DbType,omitempty"`
	DbUsername         string `xml:"DbUsername,omitempty" json:"DbUsername,omitempty"`
	DbPassword         string `xml:"DbPassword,omitempty" json:"DbPassword,omitempty"`
	DbName             string `xml:"DbName,omitempty" json:"DbName,omitempty"`
	DbHostname         string `xml:"DbHostname,omitempty" json:"DbHostname,omitempty"`
	DbPortnumber       string `xml:"DbPortnumber,omitempty" json:"DbPortnumber,omitempty"`
	FilePrefix         string `xml:"FilePrefix,omitempty" json:"FilePrefix,omitempty"`
	FileExtension      string `xml:"FileExtension,omitempty" json:"FileExtension,omitempty"`
	MigrationOutput    string `xml:"MigrationOutput,omitempty" json:"MigrationOutput,omitempty"`
	MigrationTableName string `xml:"MigrationTableName,omitempty" json:"MigrationTableName,omitempty"`
}

// Migration structure is combination of Up / Down structure.
type Migration struct {
	ID   string `xml:"id,omitempty" json:"id,omitempty"`
	Up   UpDown `xml:"up,omitempty" json:"up,omitempty"`
	Down UpDown `xml:"down,omitempty" json:"down,omitempty"`
}

type UpDown struct {
	CreateTable []CreateTable `xml:"createTable,omitempty" json:"createTable,omitempty"`
	DropTable   []DropTable   `xml:"dropTable,omitempty" json:"dropTable,omitempty"`
	AddColumn   []AddColumn   `xml:"addColumn,omitempty" json:"addColumn,omitempty"`
	DropColumn  []DropColumn  `xml:"dropColumn,omitempty" json:"dropColumn,omitempty"`
	AddIndex    []AddIndex    `xml:"addIndex,omitempty" json:"addIndex,omitempty"`
	DropIndex   []DropIndex   `xml:"dropIndex,omitempty" json:"dropIndex,omitempty"`
	Sql         string        `xml:"sql,omitempty" json:"sql,omitempty"`
}

type CreateTable struct {
	TableName string    `xml:"tableName,omitempty" json:"tableName,omitempty"`
	Columns   []Columns `xml:"columns,omitempty" json:"columns,omitempty"`
}

type DropTable struct {
	TableName string `xml:"tableName,omitempty" json:"tableName,omitempty"`
}

type AddColumn struct {
	TableName string    `xml:"tableName,omitempty" json:"tableName,omitempty"`
	Columns   []Columns `xml:"columns,omitempty" json:"columns,omitempty"`
}

type DropColumn struct {
	TableName string    `xml:"tableName,omitempty" json:"tableName,omitempty"`
	Columns   []Columns `xml:"columns,omitempty" json:"columns,omitempty"`
}

type AddIndex struct {
	TableName string    `xml:"tableName,omitempty" json:"tableName,omitempty"`
	IndexType string    `xml:"indexType,omitempty" json:"indexType,omitempty"`
	Columns   []Columns `xml:"columns,omitempty" json:"columns,omitempty"`
}

type DropIndex struct {
	TableName string    `xml:"tableName,omitempty" json:"tableName,omitempty"`
	IndexType string    `xml:"indexType,omitempty" json:"indexType,omitempty"`
	Columns   []Columns `xml:"columns,omitempty" json:"columns,omitempty"`
}

type Columns struct {
	FieldName string `xml:"fieldname,omitempty" json:"fieldname,omitempty"`
	DataType  string `xml:"datatype,omitempty" json:"datatype,omitempty"`
}
