# Artransfer CLI

```sh
> ./artransfer -h

Usage: artransfer --db="mysql" <command>

Flags:
  -h, --help                 Show context-sensitive help.
      --db="mysql"           数据库 - 类型 [可选: mysql|sqlite|postgres|sqlserver]
      --file=STRING          数据库 - 文件 [指定 SQLite 数据库文件路径]
      --name=STRING          数据库 - 名称
      --host="localhost"     数据库 - 地址
      --port=3306            数据库 - 端口
      --user="root"          数据库 - 账户
      --password=STRING      数据库 - 密码
      --table-prefix=""      数据库 - 表前缀
      --charset="utf8mb4"    数据库 - 编码
      --dsn=STRING           数据库 - DSN
  -o, --output=STRING        导出文件名

Commands:
  typecho --db="mysql"
    从 Typecho 导出 Artrans

  waline --db="mysql"
    从 Waline 导出 Artrans

Run "artransfer <command> --help" for more information on a command.
```

## 举个栗子

### Typecho

从 Typecho 数据库中导出 Artrans 文件，执行：

```sh
./artransfer typecho \
    --db="mysql"
    --host="localhost" \
    --port="3306" \
    --user="root" \
    --password="123456" \
    --name="typecho_数据库名"
```

一份 Artrans 格式的评论数据文件将导出到当前目录：

```sh
> ls
typecho-20220424-202246.artrans
```

查看更多参数可执行：`./artransfer typecho -h`

### Waline

从 Waline 数据库中导出 Artrans 文件，执行：

```sh
./artransfer waline \
    --db="mysql" \
    --host="localhost" \
    --port="3306" \
    --user="root" \
    --password="123456" \
    --name="waline_数据库名" \
    --table-prefix="wl_"
```
