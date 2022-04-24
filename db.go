package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/ArtalkJS/Artransfer-CLI/lib"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func OpenDB(ctx *lib.Context) (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   ctx.DBConf.TablePrefix,
			SingularTable: true,
			NoLowerCase:   true,
		},
	}

	dsn := ctx.DBConf.Dsn

	switch strings.ToLower(ctx.DBConf.Type) {

	case "sqlite", "sqlite3":
		if dsn == "" {
			if ctx.DBConf.File == "" {
				log.Fatal("请使用参数 --file 指定一个 SQLite 数据库路径")
			}

			dsn = ctx.DBConf.File
		}

		return gorm.Open(sqlite.Open(dsn), gormConfig)

	case "postgres", "pgsql":
		if dsn == "" {
			dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
				ctx.DBConf.Host,
				ctx.DBConf.User,
				ctx.DBConf.Password,
				ctx.DBConf.Name,
				ctx.DBConf.Port)
		}

		return gorm.Open(postgres.Open(dsn), gormConfig)

	case "mysql", "mssql", "sqlserver":
		if dsn == "" {
			dsn = fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
				ctx.DBConf.User,
				ctx.DBConf.Password,
				ctx.DBConf.Host,
				ctx.DBConf.Port,
				ctx.DBConf.Name,
				ctx.DBConf.Charset,
			)
		}

		if strings.EqualFold(ctx.DBConf.Type, "mysql") {
			return gorm.Open(mysql.Open(dsn), gormConfig)
		} else {
			return gorm.Open(sqlserver.Open(dsn), gormConfig)
		}
	}

	return nil, errors.New(`不支持的数据库类型 "` + ctx.DBConf.Type + `"`)
}
