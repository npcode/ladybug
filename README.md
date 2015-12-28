# Project Ladybug 

The simple and straightforward testcase management tools.

## Description

Project Ladubug can 

* support dashboard
* manage test case
* manage builds
* manage requirements(soon)
* support reports(soon)

## I need any kind of help! 

Since I'm new to Go language, not familiar with code convention, documentation, making excellemt code of Go language. Good at HTML/CSS/Javascript? Not at all! I don't have any chance to join web project since I start to work. So the code and poor design may(or shall?) disappoint you. Do not stay. Please make an issue or fork this repo, pull request. Every issues and pull requests are always welcome.

## Getting Started

### Prerequirements

The Project Ladybug uses below... 

* [revel](https://github.com/revel/revel) web framework.
* [gorm](https://github.com/jinzhu/gorm) database driver. 

### Installation

* You need set up database before running ladybug
* Now only Postgresql is required. You can use not only Postgresql but other relational database, but not tested. 
* Default database name is "ladybug" and user "ladybug". see "gorm.go" in app/controllers/gorm.go

### Databases

This app uses now only Postgresql. Various databases(MySQL, MarinaDB ....) will be supported. 

### Features Next

* Requirements management
* Reports
* Test environment
* Milestone
* Test coverage