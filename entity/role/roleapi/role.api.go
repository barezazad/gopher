package roleapi

import (
	"gopher/entity/access/accessenum"
	"gopher/entity/role/roleenum"
	"gopher/entity/role/rolemodel"
	"gopher/entity/service"
	"gopher/internal/core"
	"gopher/internal/core/terms"
	"gopher/internal/response"
	"net/http"

	"gopher/pkg/dictionary"

	"github.com/gin-gonic/gin"
)

// RoleAPI for injecting role service
type RoleAPI struct {
	Service service.RoleService
	Engine  *core.Engine
}

// ProvideRoleAPI for role is used in wire
func ProvideRoleAPI(c service.RoleService) RoleAPI {
	return RoleAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch a role by it's id
func (p *RoleAPI) FindByID(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, rolemodel.Table)
	var role rolemodel.Role
	var err error

	if params.ID, err = resp.GetID(c.Param("roleID"), "E1000063"); err != nil {
		return
	}

	if role, err = p.Service.FindByID(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(roleenum.ViewRole)
	resp.Status(http.StatusOK).
		MessageT(terms.VInfo, terms.Role).
		JSON(role)
}

// FindAll it return list of all role
func (p *RoleAPI) FindAll(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, rolemodel.Table)
	var roles []rolemodel.Role
	var err error

	if roles, err = p.Service.FindAll(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(roleenum.AllRole)
	resp.Status(http.StatusOK).
		MessageT(terms.AllV, terms.Roles).
		JSON(roles)
}

// return all resources
func (p *RoleAPI) Resources(c *gin.Context) {
	resp, _ := response.NewParam(p.Engine, c, rolemodel.Table)

	resp.Record(roleenum.Resources)
	resp.Status(http.StatusOK).
		MessageT(terms.AllV, terms.Resources).
		JSON(accessenum.Resources)
}

// List of roles
func (p *RoleAPI) List(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, rolemodel.Table)
	data := make(map[string]interface{})
	var err error

	if data["list"], data["count"], err = p.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(roleenum.ListRole)
	resp.Status(http.StatusOK).
		MessageT(terms.ListOfV, terms.Roles).
		JSON(data)
}

// Create role
func (p *RoleAPI) Create(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, rolemodel.Table)
	var role, createdRole rolemodel.Role
	var err error

	if err = resp.Bind(&role, "E1000064"); err != nil {
		return
	}

	role.CreatedBy = params.UserID
	if createdRole, err = p.Service.Create(role, params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.RecordCreate(roleenum.CreateRole, createdRole)
	resp.Status(http.StatusOK).
		MessageT(terms.VCreatedSuccessfully, dictionary.Translate(terms.Role)).
		JSON(createdRole)
}

// Update role
func (p *RoleAPI) Update(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, rolemodel.Table)
	var role, roleBefore, roleUpdated rolemodel.Role
	var err error

	if params.ID, err = resp.GetID(c.Param("roleID"), "E1000065"); err != nil {
		return
	}

	if err = resp.Bind(&role, "E1000066"); err != nil {
		return
	}

	role.ID = params.ID
	role.UpdatedBy = params.UserID
	if roleUpdated, roleBefore, err = p.Service.Save(role, params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(roleenum.UpdateRole, roleBefore, roleUpdated)
	resp.Status(http.StatusOK).
		MessageT(terms.VUpdatedSuccessfully, dictionary.Translate(terms.Role)).
		JSON(roleUpdated)
}

// Delete role
func (p *RoleAPI) Delete(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, rolemodel.Table)
	var role rolemodel.Role
	var err error

	if params.ID, err = resp.GetID(c.Param("roleID"), "E1000067"); err != nil {
		return
	}

	if role, err = p.Service.Delete(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(roleenum.DeleteRole, role)
	resp.Status(http.StatusOK).
		MessageT(terms.VDeletedSuccessfully, terms.Role).
		JSON()
}
