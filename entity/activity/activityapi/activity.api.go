package activityapi

import (
	"fmt"
	"gopher/entity/activity/activityenum"
	"gopher/entity/activity/activitymodel"
	"gopher/entity/service"
	"gopher/entity/user/userenum"
	"gopher/internal/core"
	"gopher/internal/core/terms"
	"gopher/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ActivityAPI for injecting activity service
type ActivityAPI struct {
	Service service.ActivityService
	Engine  *core.Engine
}

// ProvideActivityAPI for activity is used in wire
func ProvideActivityAPI(c service.ActivityService) ActivityAPI {
	return ActivityAPI{Service: c, Engine: c.Engine}
}

// Create activity
func (p *ActivityAPI) Create(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, userenum.Entity)
	var activity activitymodel.Activity

	if err := c.ShouldBindJSON(&activity); err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, err)
		return
	}

	createdActivity, err := p.Service.Save(activity, params)
	if err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	resp.Status(203).
		Message("activity created successfully").
		JSON(createdActivity)
}

// ListAll of all activities among all companies
func (p *ActivityAPI) ListAll(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, activityenum.Entity)

	data, err := p.Service.List(params)
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(activityenum.AllActivity)
	resp.Status(http.StatusOK).
		MessageT(terms.ListOfV, terms.Activities).
		JSON(data)
}

// ListSelf of all activities among all companies
func (p *ActivityAPI) ListSelf(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, activityenum.Entity)
	var err error

	params.ForceCondition = fmt.Sprintf("activities.user_id = %v", params.UserID)

	data, err := p.Service.List(params)
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(activityenum.SelfActivity)
	resp.Status(http.StatusOK).
		MessageT(terms.ListOfV, terms.Activities).
		JSON(data)
}
