Pravasan
========
Simple Migration tool intend to be used for any languages, for any db.

*Please feel free to criticize, comment, etc.*

* Install
* Usage

Install
-------

Usage
-----

Flags
```
  -d="": specify the database name
  -h="localhost": specify the database hostname
  -p=false: specify the option asking for database password
  -port="5432": specify the database port
  -prefix="": specify the database port
  -u="": specify the database username
  -version=false: print Pravasan version
```

Assuming the pravasan.conf file is set already
```
pravasan add create_table test123 id:int name:string order:int status:bool
pravasan add add_column test123 id:int
pravasan add drop_column test123 id

pravasan up
```

Work in progress are:
----
* Reading Conf file
* Creating Conf file
```
pravasan add sql 
pravasan add rename_table old_test123 new_test123
pravasan add add_index test123 id name
pravasan down [-1]
```

Few Notes: 
* moved from https://github.com/kishorevaishnav/godbmig