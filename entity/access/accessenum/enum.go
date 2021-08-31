package accessenum

const (
	Entity = "access"

	SuperAccess      = "supper:access"
	ReadDeleted      = "deleted:read"
	UserWrite        = "user:write"
	UserRead         = "user:read"
	UserAll          = "user:all"
	UserExcel        = "user:excel"
	RoleRead         = "role:read"
	RoleWrite        = "role:write"
	RoleAll          = "role:all"
	SettingRead      = "setting:read"
	SettingWrite     = "setting:write"
	ActivitySelf     = "activity:self"
	ActivityAll      = "activity:all"
	DocumentDownload = "document:download"
	DocumentDelete   = "document:delete"
	CityRead         = "city:read"
	CityWrite        = "city:write"
	CityAll          = "city:all"
	GiftRead         = "gift:read"
	GiftWrite        = "gift:write"
)

// list of resources
var Resources = []string{
	UserWrite,
	UserRead,
	UserAll,
	UserExcel,
	RoleRead,
	RoleWrite,
	RoleAll,
	SettingRead,
	SettingWrite,
	ActivitySelf,
	ActivityAll,
	DocumentDownload,
	DocumentDelete,
	CityRead,
	CityWrite,
	CityAll,
	GiftRead,
	GiftWrite,
}
