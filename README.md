#Pravasan
Simple Migration tool intend to be used for any languages, for any db.

*Please feel free to criticize, comment, etc.*
*Currently this is working for MySQL.* Soon will be available for other Databases too.

##Definition in Hindi
----
प्रवसन {pravasan} = MIGRATION(Noun)

* [Install](#install)
* [Usage](#usage)
* [High Level Features](#high-level-features)
* [All Features / Bugs](#all-features--bugs)

##Install
(in progressssss......)

##Usage

###Syntax
```pravasan [<flags>] <action> <sub-action> [data input]```

###Flags
```
Usage of pravasan:
  -confOutput="json": config file format: json, xml
  -d="": database name
  -dbType="mysql": database type
  -h="localhost": database hostname
  -indexPrefix="idx": prefix for creating Indexes
  -indexSuffix="": suffix for creating Indexes
  -migFileExtn="prvsn": migration file extension
  -migFilePrefix="": prefix for migration file
  -migOutput="json": current supported format: json, xml
  -migTableName="schema_migrations": migration table name
  -p=false: database password
  -port="5432": database port
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
pravasan add add_index test123 id order name
pravasan add create_table test123 id:int name:string order:int status:bool
pravasan add drop_column test123 id
pravasan add sql               # to add SQL statements directly.

pravasan down [-1]
pravasan up
```

If you like not to store the credentials in file then use it like this
```
pravasan -u=root -p -dbType=postgres -d=testdb -h=localhost -port=5433 add add_column test123 id:int
```

##Work in progress are:
```
pravasan add rename_table old_test123 new_test123
```
* Support for SQLite, Oracle, MongoDB, etc.,

##High Level Features
- [x] Create & read from Conf file (XML / JSON)
- [x] Output in XML, JSON format
- [x] Support for Direct SQL Statements 
- [x] Support for MySQL, Postgres

##All Features / Bugs
- [x] [v0.1](https://github.com/pravasan/pravasan/milestones/v0.1)
- [ ] [v0.2](https://github.com/pravasan/pravasan/milestones/v0.2)
- [ ] [v0.3](https://github.com/pravasan/pravasan/milestones/v0.3)
- [ ] [v1.0](https://github.com/pravasan/pravasan/milestones/v1.0)
- [ ] [v2.0](https://github.com/pravasan/pravasan/milestones/v2.0)

####Few Other Notes: 
* moved from https://github.com/kishorevaishnav/godbmig