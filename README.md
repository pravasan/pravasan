# Pravasan
Simple Migration tool intend to be used for any languages, for any db.

[![Build Status](https://travis-ci.org/pravasan/pravasan.svg?branch=master)](https://travis-ci.org/pravasan/pravasan)

*Please feel free to criticize, comment, etc.*
*Currently this is working for MySQL, PostgreSQL.* Soon will be available for other Databases too.

## Definition in Hindi
----
प्रवसन {pravasan} = MIGRATION(Noun)

* [Install](#install)
* [Usage](#usage)
* [High Level Features](#high-level-features)
* [All Features / Bugs](#all-features--bugs)
* [Details](#details)

## Install

1. Choose proper OS & Download from http://pravasan.github.io/pravasan/#download
2. Unzip / Untar the file downloaded 
3. For some OS its default but some it needs to be explicit add 
```chmod +x pravasan_*```
4. Look at the below Usage and start using it from the folder where you would like to execute & store migration files.

## Usage

### Syntax
```pravasan [<flags>] <action> <sub-action> [data input]```

### Flags
```
Usage of pravasan:
  -autoAddColumns="": Add default columns when table is created
  -confOutput="json": config file format: json, xml
  -d="": database name
  -dbType="mysql": database type
  -h="localhost": database hostname
  -indexPrefix="idx": prefix for creating Indexes
  -indexSuffix="": suffix for creating Indexes
  -migDir="./": migration file stored directory
  -migFileExtn="prvsn": migration file extension
  -migFilePrefix="": prefix for migration file
  -migOutput="json": current supported format: json, xml
  -migTableName="schema_migrations": migration table name
  -p=false: database password
  -port="5432": database port
  -storeDirectSQL=true: Store SQL in migration file instead of XML / JSON
  -u="": database username
  -version=false: print Pravasan version
```

To create the configuration file use either of the below commands & pravasan.conf.json / pravasan.conf.xml will be created
```
pravasan -u="root" -p -dbType="mysql" -d="testdb" -h="localhost" -port="5433" create conf 
pravasan -u=root -p -dbType=postgres -d=testdb -h=localhost -port=5433 -output=xml create conf 
```

Assuming the pravasan.conf.json or pravasan.conf.xml file is set already
```
pravasan add add_column test123 id:int
pravasan a ac test123 id:int                                     # Same as above

pravasan add add_index test123 id order name

pravasan add create_table test123 id:int name:string order:int status:bool
pravasan a ct test123 id:int name:string order:int status:bool              # Same as above

pravasan add drop_column test123 id
pravasan add drop_index test123 id order name
pravasan add rename_table test123 new_test123
pravasan add sql                                           # to add SQL statements directly.

pravasan down [-1]
pravasan up
pravasan up 20150103174227, 20150103174333
```

If you like not to store the credentials in file then use it like this
```
pravasan -u=root -p -dbType=postgres -d=testdb -h=localhost -port=5433 up 20150103174227
```

## Work in progress are:
- [ ] Support for Oracle, MongoDB, etc.,

## High Level Features
- [x] Create & read from Conf file (XML / JSON)
- [x] Output in XML, JSON format
- [x] Support for Direct SQL Statements 
- [x] Support for MySQL, Postgres
- Temporarily SQLite3 is stopped - due to compile issue.

## All Features / Bugs
- [x] [v0.1](https://github.com/pravasan/pravasan/milestones/v0.1)
- [x] [v0.2](https://github.com/pravasan/pravasan/milestones/v0.2)
- [x] [v0.3](https://github.com/pravasan/pravasan/milestones/v0.3)
- [ ] [v0.4](https://github.com/pravasan/pravasan/milestones/v0.4)
- [ ] [v1.0](https://github.com/pravasan/pravasan/milestones/v1.0)
- [ ] [v2.0](https://github.com/pravasan/pravasan/milestones/v2.0)

#### Few Other Notes: 
* moved from https://github.com/kishorevaishnav/godbmig

## Details

### Add few columns by default during creation of a table.
```
$ pravasan -autoAddColumns: id:int created_at:datetime modified_at:datetime create conf
```
The above command will add the 3 columns ```id```, ```created_at```, ```modified_at``` will be stored in configuration & will get added during any ```create_table``` statements.

### To store direct sql in JSON or XML format
```
$ pravasan -storeDirectSQL a ct test123 id:int
```
The above command will add (```a```) a migration & will generate create_table (```ct```) action with the table name ```test123``` with only one column ```id``` of datatype ```int``` as SQL into the JSON (by default) or XML if there is an configuration file present and set XML to generate.

### To get the list of datatypes supported

```
$ pravasan list datatypes           # All databases
$ pravasan l dt                     # All databases
$ pravasan list dt postgresql       # Will list datatypes supported by PostgreSQL
$ pravasan list datatypes sqlite3   # Will list datatypes supported by SQLite3
$ pravasan l dt mysql               # Will list datatypes supported by MySQL
```
Below is the sample output for the above last command.
```
+-------------------+-----------------------+--------------+
|     DATATYPE      | SUPPORTED DATABASE(S) |    ALIAS     |
+-------------------+-----------------------+--------------+
| VARCHAR           | MySQL                 |              |
| VARBINARY         | MySQL                 |              |
| MEDIUMBLOB        | MySQL                 |              |
| LONGBLOB          | MySQL                 |              |
| DATE              | MySQL                 |              |
| REAL              | MySQL                 |              |
| DEC               | MySQL                 | DECIMAL      |
| NUMERIC           | MySQL                 | DECIMAL      |
| STRING            | MySQL                 | VARCHAR(255) |
```