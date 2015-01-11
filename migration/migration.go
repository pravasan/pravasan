package migration

// Config used for general configuration settings for Migration application
type Config struct {
	DbHostname             string `xml:"DbHostname,omitempty" json:"DbHostname,omitempty"`
	DbName                 string `xml:"DbName,omitempty" json:"DbName,omitempty"`
	DbPassword             string `xml:"DbPassword,omitempty" json:"DbPassword,omitempty"`
	DbPort                 string `xml:"DbPort,omitempty" json:"DbPort,omitempty"`
	DbType                 string `xml:"DbType,omitempty" json:"DbType,omitempty"`
	DbUsername             string `xml:"DbUsername,omitempty" json:"DbUsername,omitempty"`
	IndexPrefix            string `xml:"IndexPrefix,omitempty" json:"IndexPrefix,omitempty"`
	IndexSuffix            string `xml:"IndexSuffix,omitempty" json:"IndexSuffix,omitempty"`
	MigrationFileExtension string `xml:"MigrationFileExtension,omitempty" json:"MigrationFileExtension,omitempty"`
	MigrationFilePrefix    string `xml:"MigrationFilePrefix,omitempty" json:"MigrationFilePrefix,omitempty"`
	MigrationOutputFormat  string `xml:"MigrationOutputFormat,omitempty" json:"MigrationOutputFormat,omitempty"`
	MigrationTableName     string `xml:"MigrationTableName,omitempty" json:"MigrationTableName,omitempty"`
}

// Migration structure is combination of Up / Down structure.
type Migration struct {
	ID   string `xml:"id,omitempty" json:"id,omitempty"`
	Up   UpDown `xml:"up,omitempty" json:"up,omitempty"`
	Down UpDown `xml:"down,omitempty" json:"down,omitempty"`
}

// UpDown structre ... #TODO need to write some proper comment
type UpDown struct {
	AddColumn   []AddColumn   `xml:"addColumn,omitempty" json:"addColumn,omitempty"`
	AddIndex    []AddIndex    `xml:"addIndex,omitempty" json:"addIndex,omitempty"`
	CreateTable []CreateTable `xml:"createTable,omitempty" json:"createTable,omitempty"`
	DropColumn  []DropColumn  `xml:"dropColumn,omitempty" json:"dropColumn,omitempty"`
	DropIndex   []DropIndex   `xml:"dropIndex,omitempty" json:"dropIndex,omitempty"`
	DropTable   []DropTable   `xml:"dropTable,omitempty" json:"dropTable,omitempty"`
	RenameTable []RenameTable `xml:"renameTable,omitempty" json:"renameTable,omitempty"`
	Sql         string        `xml:"sql,omitempty" json:"sql,omitempty"`
}

// AddColumn structre ... #TODO need to write some proper comment
type AddColumn struct {
	TableName string    `xml:"tableName,omitempty" json:"tableName,omitempty"`
	Columns   []Columns `xml:"columns,omitempty" json:"columns,omitempty"`
}

// AddIndex structre ... #TODO need to write some proper comment
type AddIndex struct {
	TableName string    `xml:"tableName,omitempty" json:"tableName,omitempty"`
	IndexType string    `xml:"indexType,omitempty" json:"indexType,omitempty"`
	Columns   []Columns `xml:"columns,omitempty" json:"columns,omitempty"`
}

// CreateTable structre ... #TODO need to write some proper comment
type CreateTable struct {
	TableName string    `xml:"tableName,omitempty" json:"tableName,omitempty"`
	Columns   []Columns `xml:"columns,omitempty" json:"columns,omitempty"`
}

// DropColumn structre ... #TODO need to write some proper comment
type DropColumn struct {
	TableName string    `xml:"tableName,omitempty" json:"tableName,omitempty"`
	Columns   []Columns `xml:"columns,omitempty" json:"columns,omitempty"`
}

// DropIndex structre ... #TODO need to write some proper comment
type DropIndex struct {
	TableName string    `xml:"tableName,omitempty" json:"tableName,omitempty"`
	IndexType string    `xml:"indexType,omitempty" json:"indexType,omitempty"`
	Columns   []Columns `xml:"columns,omitempty" json:"columns,omitempty"`
}

// DropTable structre ... #TODO need to write some proper comment
type DropTable struct {
	TableName string `xml:"tableName,omitempty" json:"tableName,omitempty"`
}

// RenameTable structre ... #TODO need to write some proper comment
type RenameTable struct {
	OldTableName string `xml:"oldTableName,omitempty" json:"oldTableName,omitempty"`
	NewTableName string `xml:"newTableName,omitempty" json:"newTableName,omitempty"`
}

// Columns structre ... #TODO need to write some proper comment
type Columns struct {
	FieldName string `xml:"fieldname,omitempty" json:"fieldname,omitempty"`
	DataType  string `xml:"datatype,omitempty" json:"datatype,omitempty"`
}
