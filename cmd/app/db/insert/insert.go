package insert

import (
	"gopher/cmd/app/db/insert/table"
	"gopher/internal/core"
)

// Insert is used for add static rows to database
func Insert(engine *core.Engine) {
	table.InsertRoles(engine)
	table.InsertUsers(engine)
	table.InsertSettings(engine)
}
