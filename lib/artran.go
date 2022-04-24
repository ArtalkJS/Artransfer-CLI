package lib

// 数据行囊 (n 个 Artran 组成一个 Artrans)
// Fields All String type (FAS)
// @link https://github.com/ArtalkJS/ArtalkGo/blob/master/model/artran.go
type Artran struct {
	ID  string `json:"id"`
	Rid string `json:"rid"`

	Content string `json:"content"`

	UA          string `json:"ua"`
	IP          string `json:"ip"`
	IsCollapsed string `json:"is_collapsed,omitempty"` // bool => string "true" or "false"
	IsPending   string `json:"is_pending,omitempty"`   // bool
	IsPinned    string `json:"is_pinned,omitempty"`    // bool

	// vote
	VoteUp   string `json:"vote_up,omitempty"`
	VoteDown string `json:"vote_down,omitempty"`

	// date
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at,omitempty"`

	// user
	Nick       string `json:"nick"`
	Email      string `json:"email"`
	Link       string `json:"link"`
	Password   string `json:"password,omitempty"`
	BadgeName  string `json:"badge_name,omitempty"`
	BadgeColor string `json:"badge_color,omitempty"`

	// page
	PageKey       string `json:"page_key"`
	PageTitle     string `json:"page_title,omitempty"`
	PageAdminOnly string `json:"page_admin_only,omitempty"` // bool

	// site
	SiteName string `json:"site_name,omitempty"`
	SiteUrls string `json:"site_urls,omitempty"`
}
