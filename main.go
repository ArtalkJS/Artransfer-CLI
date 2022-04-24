package main

import (
	"fmt"
	"strings"

	"github.com/ArtalkJS/Artransfer-CLI/lib"
	"github.com/ArtalkJS/Artransfer-CLI/typecho"
	"github.com/ArtalkJS/Artransfer-CLI/waline"
	"github.com/alecthomas/kong"

	log "github.com/sirupsen/logrus"
)

var cli struct {
	Typecho typecho.TypechoCmd `cmd:"" help:"从 Typecho 导出 Artrans"`
	Waline  waline.WalineCmd   `cmd:"" help:"从 Waline 导出 Artrans"`

	DB string `help:"数据库 - 类型 [可选: mysql|sqlite|postgres|sqlserver]" required:"" default:"mysql"`

	File string `help:"数据库 - 文件 [指定 SQLite 数据库文件路径]" type:"existingfile"`
	Name string `help:"数据库 - 名称"`

	Host     string `help:"数据库 - 地址" default:"localhost"`
	Port     int    `help:"数据库 - 端口" default:"3306"`
	User     string `help:"数据库 - 账户" default:"root"`
	Password string `help:"数据库 - 密码"`

	TablePrefix string `help:"数据库 - 表前缀" default:""`
	Charset     string `help:"数据库 - 编码" default:"utf8mb4"`
	Dsn         string `help:"数据库 - DSN"`

	Output string `help:"导出文件名" short:"o"`
}

func main() {
	kong := kong.Parse(&cli)
	ctx := &lib.Context{
		DBConf: lib.DBConf{
			Type:        cli.DB,
			File:        cli.File,
			Name:        cli.Name,
			Host:        cli.Host,
			Port:        cli.Port,
			User:        cli.User,
			Password:    cli.Password,
			TablePrefix: cli.TablePrefix,
			Charset:     cli.Charset,
			Dsn:         cli.Dsn,
		},
		OutputFile: cli.Output,
	}

	PrintDBConf(ctx)

	if db, err := OpenDB(ctx); err == nil {
		ctx.DB = db
	} else {
		log.Fatal("数据库连接失败 ", err)
	}

	err := kong.Run(ctx)
	kong.FatalIfErrorf(err)
}

func PrintDBConf(ctx *lib.Context) {
	fmt.Println()
	fmt.Println("--------------------")
	fmt.Println("     数据库配置      ")
	fmt.Println("--------------------")

	fmt.Println("类型：" + ctx.DBConf.Type)
	if ctx.DBConf.File != "" {
		fmt.Println("文件：" + ctx.DBConf.File)
	} else {
		fmt.Println("名称：" + ctx.DBConf.Name)
		fmt.Println("地址：" + ctx.DBConf.Host + fmt.Sprintf(":%d", ctx.DBConf.Port))
		fmt.Println("账户：" + ctx.DBConf.User)
		fmt.Println("密码：" + strings.Repeat("*", len([]rune(ctx.DBConf.Password))))

		displayTablePrefix := ctx.DBConf.TablePrefix
		if ctx.DBConf.TablePrefix == "" {
			displayTablePrefix = "(未设定)"
		}
		fmt.Println("表前缀：" + displayTablePrefix)
	}

	fmt.Println()
}
