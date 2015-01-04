Pravasan
========
Simple Migration tool intend to be used for any languages, for any db.

*Please feel free to criticize, comment, etc.*
*Currently this is working for MySQL.* Soon will be available for other Databases too.

Definition in Hindi
----
प्रवसन {pravasan} = MIGRATION(Noun)

* Install
* Usage

Install
-------

Usage
-----

Flags
```
Usage of pravasan:
  -d="": specify the database name
  -db_type="mysql": specify the database type
  -extn="prvsn": specify the migration file extension
  -h="localhost": specify the database hostname
  -p=false: specify the option asking for database password
  -port="5432": specify the database port
  -prefix="": specify the text to be prefix with the migration file
  -u="": specify the database username
  -version=false: print Pravasan version
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

Planning for v0.5
----

Planning for v1.0
----
Transactional Support

Few Notes: 
* moved from https://github.com/kishorevaishnav/godbmig