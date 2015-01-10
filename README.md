Pravasan
========
Simple Migration tool intend to be used for any languages, for any db.

*Please feel free to criticize, comment, etc.*
*Currently this is working for MySQL.* Soon will be available for other Databases too.

Definition in Hindi
----
प्रवसन {pravasan} = MIGRATION(Noun)

* [Install](#install)
* [Usage](#usage)
* [High Level Features](#high-level-features)
* [All Features / Bugs](#all-features--bugs)

Install
-------
(in progressssss......)

Usage
-----

Flags
```
Usage of pravasan:
  -d="": specify the database name
  -dbType="mysql": specify the database type
  -extn="prvsn": specify the migration file extension
  -h="localhost": specify the database hostname
  -migration_table_name="schema_migrations": supported format are json, xml
  -output="json": supported format are json, xml
  -p=false: specify the option asking for database password
  -port="5432": specify the database port
  -prefix="": specify the text to be prefix with the migration file
  -u="": specify the database username
  -version=false: print Pravasan version
```

To create the configuration file use either of the below commands & pravasan.conf.json / pravasan.conf.xml will be created
```
pravasan -u="root" -p -d="testdb" -h="localhost" -port="5433" create conf 
pravasan -u="root" -p -d="testdb" -h="localhost" -port="5433" -output="xml" create conf 
```

Assuming the pravasan.conf.json file is set already
```
pravasan add create_table test123 id:int name:string order:int status:bool
pravasan add add_column test123 id:int
pravasan add drop_column test123 id

pravasan up
pravasan down [-1]
```

Work in progress are:
----
```
pravasan add sql 
pravasan add rename_table old_test123 new_test123
pravasan add add_index test123 id name
```
* Fix _ in field names 
* Creating Conf file
* Support for Postgres, SQLite, Oracle, MongoDB, etc.,

High Level Features
----
- [x] Output in XML, JSON format
- [x] Support for MySQL
- [x] Create & read from Conf file (XML / JSON)

All Features / Bugs
========
- [ ] [v0.1](https://github.com/pravasan/pravasan/milestones/v0.1)
- [ ] [v0.2](https://github.com/pravasan/pravasan/milestones/v0.2)
- [ ] [v0.3](https://github.com/pravasan/pravasan/milestones/v0.3)
- [ ] [v1.0](https://github.com/pravasan/pravasan/milestones/v1.0)
- [ ] [v2.0](https://github.com/pravasan/pravasan/milestones/v2.0)

Few Notes: 
* moved from https://github.com/kishorevaishnav/godbmig