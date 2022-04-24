package typecho

// Typecho 评论数据表
type TypechoComment struct {
	Coid     int    `gorm:"column:coid"` // comment_id
	Cid      int    `gorm:"column:cid"`  // content_id
	Created  int    `gorm:"column:created"`
	Author   string `gorm:"column:author"`
	AuthorId int    `gorm:"column:authorId"`
	OwnerId  int    `gorm:"column:ownerId"`
	Mail     string `gorm:"column:mail"`
	Url      string `gorm:"column:url"`
	Ip       string `gorm:"column:ip"`
	Agent    string `gorm:"column:agent"`
	Text     string `gorm:"column:text"`
	Type     string `gorm:"column:type"`
	Status   string `gorm:"column:status"`
	Parent   int    `gorm:"column:parent"`
	Stars    int    `gorm:"column:stars"`
	Notify   int    `gorm:"column:notify"`
	Likes    int    `gorm:"column:likes"`
	Dislikes int    `gorm:"column:dislikes"`
}

// Typecho 内容数据表 (Type: post, page)
type TypechoContent struct {
	Cid          int    `gorm:"column:cid"`
	Title        string `gorm:"column:title"`
	Slug         string `gorm:"column:slug"`
	Created      int    `gorm:"column:created"`
	Modified     int    `gorm:"column:modified"`
	Text         string `gorm:"column:text"`
	Order        int    `gorm:"column:order"`
	AuthorId     int    `gorm:"column:authorId"`
	Template     string `gorm:"column:template"`
	Type         string `gorm:"column:type"`
	Status       string `gorm:"column:status"`
	Password     string `gorm:"column:password"`
	CommentsNum  int    `gorm:"column:commentsNum"`
	AllowComment string `gorm:"column:allowComment"`
	AllowPing    string `gorm:"column:allowPing"`
	AllowFeed    string `gorm:"column:allowFeed"`
	Parent       int    `gorm:"column:parent"`
	Views        int    `gorm:"column:views"`
	ViewsNum     int    `gorm:"column:viewsNum"`
	LikesNum     int    `gorm:"column:likesNum"`
	WordCount    int    `gorm:"column:wordCount"`
	Likes        int    `gorm:"column:likes"`
}

// Typecho 配置数据表
type TypechoOption struct {
	Name  string `gorm:"column:name"`
	User  int    `gorm:"column:user"`
	Value string `gorm:"column:value"`
}

// Typecho 关联表 (Content => Metas)
type TypechoRelationship struct {
	Cid int `gorm:"column:cid"` // content_id
	Mid int `gorm:"column:mid"` // meta_id
}

// Typecho 附加字段数据表
type TypechoMeta struct {
	Mid         int    `gorm:"column:mid"`
	Name        string `gorm:"column:name"`
	Slug        string `gorm:"column:slug"`
	Type        string `gorm:"column:type"`
	Description string `gorm:"column:description"`
}
