package migration

type Config struct {
	Db_type          string `xml:"Db_type,omitempty"json:"Db_type,omitempty"`
	Db_username      string `xml:"Db_username,omitempty"json:"Db_username,omitempty"`
	Db_password      string `xml:"Db_password,omitempty"json:"Db_password,omitempty"`
	Db_name          string `xml:"Db_name,omitempty"json:"Db_name,omitempty"`
	Db_hostname      string `xml:"Db_hostname,omitempty"json:"Db_hostname,omitempty"`
	Db_portnumber    string `xml:"Db_portnumber,omitempty"json:"Db_portnumber,omitempty"`
	File_prefix      string `xml:"File_prefix,omitempty"json:"File_prefix,omitempty"`
	File_extension   string `xml:"File_extension,omitempty"json:"File_extension,omitempty"`
	Migration_output string `xml:"Migration_output,omitempty"json:"Migration_output,omitempty"`
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
