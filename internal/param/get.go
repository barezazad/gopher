package param

import (
	"gopher/internal/core"
	"gopher/pkg/dictionary"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Get is a function for filling param.Model
func Get(c *gin.Context, engine *core.Engine, entity string) (param Param) {
	var err error

	generateOrder(c, &param, entity)
	generateSelectedColumns(c, &param)
	generateLimit(c, engine, &param)
	generateOffset(c, engine, &param)

	// param.Search = strings.TrimSpace(c.Query("search"))
	param.Filter = strings.TrimSpace(c.Query("filter"))

	userID, ok := c.Get("USER_ID")
	if ok {
		param.UserID = userID.(uint)
	} else {
		engine.ServerLog.CheckInfo(err, "User ID is not exist")
	}

	username, ok := c.Get("USERNAME")
	if ok {
		param.Username = username.(string)
	} else {
		engine.ServerLog.CheckInfo(err, "Username is not exist")
	}

	if c.Query("deleted") == "true" {
		param.ShowDeletedRows = true
	}

	lang, ok := c.Get("LANGUAGE")
	if ok {
		param.Lang = lang.(string)
	} else {
		param.Lang = dictionary.En
	}

	return param
}

func generateOrder(c *gin.Context, param *Param, entity string) {
	orderBy := entity + ".id"
	direction := "desc"

	if c.Query("order_by") != "" {
		orderBy = c.Query("order_by")
	}

	if c.Query("direction") != "" {
		direction = c.Query("direction")
	}

	param.Order = orderBy + " " + direction
}

func generateSelectedColumns(c *gin.Context, param *Param) {
	param.Select = "*"
	if c.Query("select") != "" {
		param.Select = c.Query("select")
	}
}

func generateLimit(c *gin.Context, engine *core.Engine, param *Param) {
	var err error
	param.Limit = 10
	if c.Query("page_size") != "" {
		param.Limit, err = strconv.Atoi(c.Query("page_size"))
		if err != nil {
			// TODO: get path from gin.Context
			engine.ServerLog.CheckError(err, "E1000164", "Limit is not a number")
			param.Limit = 10
		}
	}
}

func generateOffset(c *gin.Context, engine *core.Engine, param *Param) {
	var page int
	var err error
	if c.Query("page") != "" {
		page, err = strconv.Atoi(c.Query("page"))
		if err != nil {
			// TODO: get path from gin.Context
			engine.ServerLog.CheckError(err, "E1000165", "Offset is not a positive number")
			page = 0
		}
	}

	param.Offset = param.Limit * (page)
}
