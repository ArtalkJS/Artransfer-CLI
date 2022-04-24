package typecho

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/ArtalkJS/Artransfer-CLI/lib"
	"github.com/elliotchance/phpserialize"
	log "github.com/sirupsen/logrus"
)

const (
	// 已测试可用版本
	TypechoTestedVerMain      = "1.1"
	TypechoTestedVerSub       = "17.10.30"
	TypechoRewritePostDefault = "/archives/{cid}/"
	TypechoRewritePageDefault = "/{slug}.html"
)

type TypechoCmd struct {
	RewritePost string `help:"文章 URL 重写规则"`
	RewritePage string `help:"独立页面 URL 重写规则"`
}

// Typecho 升级相关的代码 @see https://github.com/typecho/typecho/blob/64b8e686885d8ab4c7f0cdc3d6dc2d99fa48537c/var/Utils/Upgrade.php
// 路由 @see https://github.com/typecho/typecho/blob/530312443142577509df88ce88cf3274fac9b8c4/var/Widget/Options/Permalink.php#L319
// DB @see https://github.com/typecho/typecho/blob/6558fd5e030a950335d53038f82728b06ad6c32d/install/Mysql.sql
func (cmd *TypechoCmd) Run(ctx *lib.Context) error {
	ctx.SrcType = "typecho"

	t := Typecho{
		Cmd: cmd,
	}

	if ctx.DBConf.TablePrefix == "" {
		ctx.DBConf.TablePrefix = "typecho_"
	}

	// Load Options
	tbOptions := ctx.DBConf.TablePrefix + "options"
	ctx.DB.Raw("SELECT * FROM " + tbOptions).Scan(&t.Options)
	fmt.Println(fmt.Sprintf("从数据表 `%s` 获取 %d 条记录", tbOptions, len(t.Options)))

	for _, opt := range t.Options {
		switch opt.Name {
		case "generator":
			t.SrcVersion = opt.Value
		case "title":
			t.SrcSiteName = opt.Value
		case "siteUrl":
			t.SrcSiteURL = opt.Value
		}
	}

	// 检查数据源版本号
	t.VersionCheck()

	// 重写规则
	t.RewriteRuleReady()

	// load Metas
	tbMetas := ctx.DBConf.TablePrefix + "metas"
	ctx.DB.Raw("SELECT * FROM " + tbMetas).Scan(&t.Metas)
	fmt.Println(fmt.Sprintf("从数据表 `%s` 获取 %d 条记录", tbMetas, len(t.Metas)))

	// load Relationships
	tbRelationships := ctx.DBConf.TablePrefix + "relationships"
	ctx.DB.Raw("SELECT * FROM " + tbRelationships).Scan(&t.Relationships)
	fmt.Println(fmt.Sprintf("从数据表 `%s` 获取 %d 条记录", tbRelationships, len(t.Relationships)))

	// 获取 contents
	tbContents := ctx.DBConf.TablePrefix + "contents"
	ctx.DB.Raw(fmt.Sprintf("SELECT * FROM %s ORDER BY created ASC", tbContents)).Scan(&t.Contents)
	fmt.Println(fmt.Sprintf("从数据表 `%s` 获取 %d 条记录", tbContents, len(t.Contents)))

	// 获取 comments
	tbComments := ctx.DBConf.TablePrefix + "comments"
	ctx.DB.Raw(fmt.Sprintf("SELECT * FROM %s ORDER BY created ASC", tbComments)).Scan(&t.Comments)
	fmt.Println(fmt.Sprintf("从数据表 `%s` 获取 %d 条记录", tbComments, len(t.Comments)))

	// 导出前参数汇总
	print("\n")

	print("# 请过目：\n\n")

	// 显示第一条数据
	pageKeyEgPost := ""
	pageKeyEgPage := ""
	for _, c := range t.Contents {
		if c.Type == "post" {
			firstPost := lib.SprintEncodeData("第一篇文章", t.Contents[0])
			print(lib.HideJsonLongText("Text", firstPost))
			pageKeyEgPost = t.GetNewPageKey(c)
			fmt.Printf(" -> 生成 PageKey: %#v\n\n", pageKeyEgPost)
			break
		}
	}

	for _, c := range t.Contents {
		if c.Type == "page" {
			pageKeyEgPage = t.GetNewPageKey(c)
			break
		}
	}

	if len(t.Comments) > 0 {
		lib.PrintEncodeData("第一条评论", t.Comments[0])
	}

	lib.PrintTable([][]interface{}{
		{"[基本信息]", "读取数据"},
		{"站点名称", fmt.Sprintf("%#v", t.SrcSiteName)},
		{"BaseURL", fmt.Sprintf("%#v", t.SrcSiteURL)},
		{fmt.Sprintf("评论: %d", len(t.Comments)), fmt.Sprintf("页面: %d", len(t.Contents))},
	})

	lib.PrintTable([][]interface{}{
		{"[重写规则]", "用于生成 pageKey (评论页面唯一标识)", "生成示例"},
		{"文章页面", fmt.Sprintf("%#v", cmd.RewritePost), pageKeyEgPost},
		{"独立页面", fmt.Sprintf("%#v", cmd.RewritePage), pageKeyEgPage},
	})

	print("\n")
	println("若以上内容不符合预期，请向我们反馈：https://artalk.js.org/guide/transfer.html")
	print("\n")

	// 准备导出评论
	println()

	// 开始执行导出
	ctx.Export(ctx.ArtransToJson(t.ToArtrans()))

	return nil
}

type Typecho struct {
	Cmd *TypechoCmd

	SrcVersion  string
	SrcSiteName string
	SrcSiteURL  string

	Comments      []TypechoComment
	Options       []TypechoOption
	Contents      []TypechoContent
	Metas         []TypechoMeta
	Relationships []TypechoRelationship

	OptionRoutingTable map[string]TypechoRoute
}

func (t *Typecho) ToArtrans() []lib.Artran {
	contents := t.Contents
	comments := t.Comments

	print("\n====================================\n\n")
	fmt.Println(fmt.Sprintf("[开始导出] 共 %d 个页面，%d 条评论", len(contents), len(comments)))
	print("\n")

	artrans := []lib.Artran{}

	// 导出页面
	for _, c := range contents {
		// 导出评论
		commentTotal := 0

		for _, co := range comments {
			if co.Cid != c.Cid {
				continue
			}

			artran := lib.Artran{
				ID:  lib.ToString(co.Coid),
				Rid: lib.ToString(co.Parent),

				Content:   co.Text,
				PageKey:   t.GetNewPageKey(c),
				PageTitle: c.Title,
				SiteName:  t.SrcSiteName,

				Nick:  co.Author,
				Email: co.Mail,
				Link:  co.Url,

				UA: co.Agent,
				IP: co.Ip,

				IsPending: lib.ToString(co.Status != "approved"),

				CreatedAt: fmt.Sprintf("%v", co.Created),
			}

			artrans = append(artrans, artran)

			commentTotal++
		}

		fmt.Printf("+ [%-3d] 条评论 <- [%5d] %-30s | %#v\n", commentTotal, c.Cid, c.Slug, c.Title)
	}

	return artrans
}

// 获取新的 PageKey (根据重写规则)
func (t *Typecho) GetNewPageKey(content TypechoContent) string {
	date := lib.ParseDate(fmt.Sprintf("%v", content.Created))

	// 替换内容制作
	replaces := map[string]string{
		"cid":      fmt.Sprintf("%v", content.Cid),
		"slug":     content.Slug,
		"category": t.GetContentCategory(content.Cid),
		"year":     fmt.Sprintf("%v", date.Local().Year()),
		"month":    fmt.Sprintf("%v", date.Local().Month()),
		"day":      fmt.Sprintf("%v", date.Local().Day()),
	}

	rewriteRule := t.Cmd.RewritePost
	if strings.HasPrefix(content.Type, "post") {
		rewriteRule = t.Cmd.RewritePost
	} else if strings.HasPrefix(content.Type, "page") {
		rewriteRule = t.Cmd.RewritePage
	}

	rewriteRule = "/" + strings.TrimPrefix(rewriteRule, "/")
	return t.ReplaceAllBracketed(rewriteRule, replaces)
}

// 替换字符串
// @note 特别注意：replaces 的 key 外侧无 “[]” 或 “{}”
func (t *Typecho) ReplaceAllBracketed(data string, replaces map[string]string) string {
	r := regexp.MustCompile(`(\[|\{)\s*(.*?)\s*(\]|\})`) // 同时支持 {} 和 []
	return r.ReplaceAllStringFunc(data, func(m string) string {
		key := r.FindStringSubmatch(m)[2]
		if val, isExist := replaces[key]; isExist {
			return val
		} else {
			log.Error(fmt.Sprintf("[重写规则] \"%s\" 变量无效", key))
		}
		return m
	})
}

// 版本过旧检测
func (t *Typecho) VersionCheck() {
	r := regexp.MustCompile(`Typecho[\s]*(.+)\/(.+)`)
	group := r.FindStringSubmatch(t.SrcVersion)
	if len(group) < 2 {
		log.Warn(`无法确认您的 Typecho 版本号："` + fmt.Sprintf("%#v", t.SrcVersion) + `"`)
		return
	}

	verMain := lib.ParseVersion(group[1])
	verSub := lib.ParseVersion(group[2])

	if verMain < lib.ParseVersion(TypechoTestedVerMain) ||
		verSub < lib.ParseVersion(TypechoTestedVerSub) {
		print("\n")
		log.Warn(fmt.Sprintf("Typecho 当前版本 \"%s\" 旧于 \"Typecho %s/%s\"",
			t.SrcVersion, TypechoTestedVerMain, TypechoTestedVerSub))
		log.Warn("不能保证导出数据的完整性，但你可以选择升级 Typecho: http://docs.typecho.org/upgrade")
	}
}

// TypechoRoute 路由
type TypechoRoute struct {
	URL    string `json:"url"`
	Format string `json:"format"`
	Params string `json:"params"`
	Regx   string `json:"regx"`
	Action string `json:"action"`
	Widget string `json:"widget"`
}

// 获取 typecho_options 表里面的 name:routingTable option，解析 option.value。
// option.value 数据为 PHP 序列化后的结果。
// @see https://www.php.net/manual/en/function.serialize.php
// @see https://www.php.net/manual/en/function.unserialize.php
func (t *Typecho) GetOptionRoutingTable() (map[string]TypechoRoute, error) {
	// 仅获取一次
	if t.OptionRoutingTable != nil {
		return t.OptionRoutingTable, nil
	}

	dataStr := ""
	for _, opt := range t.Options {
		if opt.Name == "routingTable" {
			dataStr = opt.Value
			break
		}
	}

	if dataStr == "" {
		return map[string]TypechoRoute{}, errors.New("`routingTable` Not Found in Options")
	}

	// Unmarshal
	var data map[interface{}]interface{}
	err := phpserialize.Unmarshal([]byte(dataStr), &data)
	if err != nil {
		return map[string]TypechoRoute{}, err
	}

	// interface{} to struct
	routingTable := map[string]TypechoRoute{}
	for k, v := range data {
		var r TypechoRoute
		json.Unmarshal([]byte(fmt.Sprintf("%v", v)), &r)
		routingTable[fmt.Sprintf("%v", k)] = r
	}

	t.OptionRoutingTable = routingTable // 就不再反复获取了

	return routingTable, nil
}

// 获取一个 Route
func (t *Typecho) GetOptionRoute(name string) (TypechoRoute, error) {
	routingTable, err := t.GetOptionRoutingTable()
	if err != nil {
		return TypechoRoute{}, err
	}

	for n, route := range routingTable {
		if n == name {
			return route, nil
		}
	}

	return TypechoRoute{}, errors.New(`Route Name "` + name + `" Not Found`)
}

// 重写路径获取
func (t *Typecho) RewriteRuleReady() {
	check := func(nameText string, routeName string, field *string, defaultVal string) {
		if *field == "" {
			// 从数据库获取
			route, err := t.GetOptionRoute(routeName)
			if err != nil || route.URL == "" {
				if err != nil {
					log.Error(err)
				}

				*field = defaultVal
				log.Error("[重写规则] \"" + nameText + "\" 无法从数据库读取，将使用默认值 \"" + t.Cmd.RewritePost + "\"")
				return
			}

			readRule := route.URL

			// 将 [year:digital:4] 变为 [year]
			r := regexp.MustCompile(`\[(([a-zA-Z0-9]+):?(.*?))\]`)
			subMatches := r.FindAllStringSubmatch(readRule, -1)
			replaces := map[string]string{}
			for _, m := range subMatches {
				if len(m) < 1 {
					continue
				}

				// 0 => [year:digital:4]
				// 1 => year:digital:4
				// 2 => year
				// 3 => digital:4
				replaces[m[1]] = "[" + m[2] + "]"
			}
			readRule = t.ReplaceAllBracketed(readRule, replaces) // replaces keys 无括号

			// 保存从数据库读取到的 rule
			*field = readRule
			fmt.Println("重写规则 \"" + nameText + "\" 读取成功")
		} else {
			fmt.Println("[重写规则] 自定义 \""+nameText+"\" 规则：", fmt.Sprintf("%#v", *field))
		}
	}

	check("文章页", "post", &t.Cmd.RewritePost, TypechoRewritePostDefault)
	check("独立页面", "page", &t.Cmd.RewritePage, TypechoRewritePageDefault)
}

// 获取 Content 的 Metas
func (t *Typecho) GetContentMetas(cid int) []TypechoMeta {
	metaIds := []int{}
	for _, rela := range t.Relationships {
		if rela.Cid == cid {
			metaIds = append(metaIds, rela.Mid)
		}
	}

	metas := []TypechoMeta{}
	for _, m := range t.Metas {
		isNeed := false
		for _, id := range metaIds {
			if id == m.Mid {
				isNeed = true
				break
			}
		}

		if isNeed {
			metas = append(metas, m)
		}
	}

	return metas
}

// 获取 Content 的分类
func (t *Typecho) GetContentCategory(cid int) string {
	metas := t.GetContentMetas(cid)
	for _, m := range metas {
		if m.Type == "category" {
			return m.Slug
		}
	}

	return ""
}
