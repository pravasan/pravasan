package migration

type Config struct {
	Db_username   string
	Db_password   string
	Db_name       string
	Db_hostname   string
	Db_portnumber string
	File_prefix   string
}

type Migration struct {
	Id   string `xml:"id,omitempty" json:"id,omitempty"`
	Up   UpDown `xml:"up,omitempty" json:"up,omitempty"`
	Down UpDown `xml:"down,omitempty" json:"down,omitempty"`
}

type UpDown struct {
	Create_Table []CreateTable `xml:"create_table,omitempty" json:"create_table,omitempty"`
	Drop_Table   []DropTable   `xml:"drop_table,omitempty" json:"drop_table,omitempty"`
	Add_Column   []AddColumn   `xml:"add_column,omitempty" json:"add_column,omitempty"`
	Drop_Column  []DropColumn  `xml:"drop_column,omitempty" json:"drop_column,omitempty"`
	Add_Index    []AddIndex    `xml:"add_index,omitempty" json:"add_index,omitempty"`
	Drop_Index   []DropIndex   `xml:"drop_index,omitempty" json:"drop_index,omitempty"`
	Sql          string        `xml:"sql,omitempty" json:"sql,omitempty"`
}

type CreateTable struct {
	Table_Name string    `xml:"table_name,omitempty" json:"table_name,omitempty"`
	Columns    []Columns `xml:"columns,omitempty" json:"columns,omitempty"`
}

type DropTable struct {
	Table_Name string `xml:"table_name,omitempty" json:"table_name,omitempty"`
}

type AddColumn struct {
	Table_Name string    `xml:"table_name,omitempty" json:"table_name,omitempty"`
	Columns    []Columns `xml:"columns,omitempty" json:"columns,omitempty"`
}

type DropColumn struct {
	Table_Name string    `xml:"table_name,omitempty" json:"table_name,omitempty"`
	Columns    []Columns `xml:"columns,omitempty" json:"columns,omitempty"`
}

type AddIndex struct {
	Table_Name string    `xml:"table_name,omitempty" json:"table_name,omitempty"`
	Index_Type string    `xml:"index_type,omitempty" json:"index_type,omitempty"`
	Columns    []Columns `xml:"columns,omitempty" json:"columns,omitempty"`
}

type DropIndex struct {
	Table_Name string    `xml:"table_name,omitempty" json:"table_name,omitempty"`
	Index_Type string    `xml:"index_type,omitempty" json:"index_type,omitempty"`
	Columns    []Columns `xml:"columns,omitempty" json:"columns,omitempty"`
}

type Columns struct {
	FieldName string `xml:"fieldname,omitempty" json:"fieldname,omitempty"`
	DataType  string `xml:"datatype,omitempty" json:"datatype,omitempty"`
}
