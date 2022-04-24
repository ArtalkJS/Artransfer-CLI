package lib

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

type Context struct {
	SrcType    string
	OutputFile string

	DB     *gorm.DB
	DBConf DBConf
}

type DBConf struct {
	Type string

	File string
	Name string

	Host     string
	Port     int
	User     string
	Password string

	TablePrefix string
	Charset     string
	Dsn         string
}

func (ctx *Context) Export(json string) {
	filename := ctx.OutputFile
	if filename == "" {
		t := time.Now()
		filename = ctx.SrcType + "-" + t.Format("20060102-150405") + ".artrans"
	}

	dst, err := os.Create(filename)
	if err != nil {
		log.Fatal("导出文件创建失败 ", err)
	}
	defer dst.Close()

	// JSON 写入文件
	if _, err = dst.Write([]byte(json)); err != nil {
		log.Fatal("导出文件写入失败 ", err)
	}

	absPath, _ := filepath.Abs(filename)

	fmt.Println("成功导出到文件 " + absPath)
}

func (ctx *Context) ArtransToJson(artrans []Artran, pretty ...bool) string {
	var err error
	var jsonBuf []byte
	if len(pretty) > 0 && pretty[0] {
		jsonBuf, err = json.MarshalIndent(artrans, "", " ")
	} else {
		jsonBuf, err = json.Marshal(artrans)
	}

	if err != nil {
		log.Fatal("JSON 输出错误 ", err)
	}

	return string(jsonBuf)
}
