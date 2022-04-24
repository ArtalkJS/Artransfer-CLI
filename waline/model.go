package waline

import "time"

// @link https://github.com/walinejs/waline/blob/main/assets/waline.sql

// TABLE `wl_Comment`
type Comment struct {
	ID         int64     `gorm:"column:id" json:"id"`
	UserId     int64     `gorm:"column:user_id" json:"user_id"`
	Comment    string    `gorm:"column:comment" json:"comment"`
	InsertedAt time.Time `gorm:"column:insertedAt" json:"insertedAt"`
	Ip         string    `gorm:"column:ip" json:"ip"`
	Link       string    `gorm:"column:link" json:"link"`
	Mail       string    `gorm:"column:mail" json:"mail"`
	Nick       string    `gorm:"column:nick" json:"nick"`
	Pid        int64     `gorm:"column:pid" json:"pid"`
	Rid        int64     `gorm:"column:rid" json:"rid"`
	Status     string    `gorm:"column:status" json:"status"`
	Ua         string    `gorm:"column:ua" json:"ua"`
	Url        string    `gorm:"column:url" json:"url"`
	CreatedAt  time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

// TABLE `wl_Counter`
type Counter struct {
	ID        int64     `gorm:"column:id" json:"id"`
	Time      int64     `gorm:"column:time" json:"time"`
	Url       string    `gorm:"column:url" json:"url"`
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

// TABLE `wl_Users`
type Users struct {
	ID          int64     `gorm:"column:id" json:"id"`
	DisplayName string    `gorm:"column:display_name" json:"display_name"`
	Email       string    `gorm:"column:email" json:"email"`
	Password    string    `gorm:"column:password" json:"password"`
	Type        string    `gorm:"column:type" json:"type"`
	Url         string    `gorm:"column:url" json:"url"`
	Avatar      string    `gorm:"column:avatar" json:"avatar"`
	Github      string    `gorm:"column:github" json:"github"`
	Twitter     string    `gorm:"column:twitter" json:"twitter"`
	Facebook    string    `gorm:"column:facebook" json:"facebook"`
	Google      string    `gorm:"column:google" json:"google"`
	Weibo       string    `gorm:"column:weibo" json:"weibo"`
	Qq          string    `gorm:"column:qq" json:"qq"`
	TwoFA       string    `gorm:"column:2fa" json:"2fa"`
	CreatedAt   time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}
