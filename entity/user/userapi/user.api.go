package userapi

import (
	"fmt"
	"gopher/entity/service"
	"gopher/entity/user/userenum"
	"gopher/entity/user/usermodel"
	"gopher/internal/core"
	"gopher/internal/core/terms"
	"gopher/internal/response"
	"gopher/pkg/generr"
	"gopher/pkg/helper/excel"
	"net/http"

	"gopher/pkg/dictionary"

	"github.com/gin-gonic/gin"
)

// UserAPI for injecting user service
type UserAPI struct {
	Service service.UserService
	Engine  *core.Engine
}

// ProvideUserAPI for user is used in wire
func ProvideUserAPI(c service.UserService) UserAPI {
	return UserAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch a user by it's id
func (p *UserAPI) FindByID(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, usermodel.Table)
	var user usermodel.User
	var err error

	if params.ID, err = resp.GetID(c.Param("userID"), "E1000142"); err != nil {
		return
	}

	if user, err = p.Service.FindByID(params); err != nil {
		resp.Error(err).JSON()
		return
	}
	user.Password = ""

	resp.Record(userenum.ViewUser)
	resp.Status(http.StatusOK).
		MessageT(terms.VInfo, terms.User).
		JSON(user)
}

// FindAll it return list of all user
func (p *UserAPI) FindAll(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, usermodel.Table)
	var users []usermodel.User
	var err error

	if users, err = p.Service.FindAll(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(userenum.AllUser)
	resp.Status(http.StatusOK).
		MessageT(terms.AllV, terms.Users).
		JSON(users)
}

// FindByUsername is used when we try to find a user with username
func (p *UserAPI) FindByUsername(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, usermodel.Table)
	var err error

	params.Search = c.Param("username")
	user, err := p.Service.FindByUsername(params)
	if err != nil {
		resp.Error(err).JSON()
		return
	}
	user.Password = ""

	resp.Record(userenum.ViewUser)
	resp.Status(http.StatusOK).
		MessageT(terms.VInfo, terms.User).
		JSON(user)
}

// List of users
func (p *UserAPI) List(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, usermodel.Table)
	data := make(map[string]interface{})
	var err error

	if data["list"], data["count"], err = p.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(userenum.ListUser)
	resp.Status(http.StatusOK).
		MessageT(terms.ListOfV, terms.Users).
		JSON(data)
}

// Create user
func (p *UserAPI) Create(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, usermodel.Table)
	var user, createdUser usermodel.User
	var err error

	if err = resp.Bind(&user, "E1000143"); err != nil {
		return
	}

	user.CreatedBy = params.UserID
	if createdUser, err = p.Service.Create(user, params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.RecordCreate(userenum.CreateUser, createdUser)
	resp.Status(http.StatusOK).
		MessageT(terms.VCreatedSuccessfully, dictionary.Translate(terms.User)).
		JSON(createdUser)
}

// Update user
func (p *UserAPI) Update(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, usermodel.Table)
	var user, userBefore, userUpdated usermodel.User
	var err error

	if params.ID, err = resp.GetID(c.Param("userID"), "E1000144"); err != nil {
		return
	}

	if err = resp.Bind(&user, "E1000145"); err != nil {
		return
	}

	user.ID = params.ID
	user.UpdatedBy = params.UserID
	if userUpdated, userBefore, err = p.Service.Save(user, params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(userenum.UpdateUser, userBefore, userUpdated)
	resp.Status(http.StatusOK).
		MessageT(terms.VUpdatedSuccessfully, dictionary.Translate(terms.User)).
		JSON(userUpdated)
}

// Delete user
func (p *UserAPI) Delete(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, usermodel.Table)
	var user usermodel.User
	var err error

	if params.ID, err = resp.GetID(c.Param("userID"), "E1000146"); err != nil {
		return
	}

	if user, err = p.Service.Delete(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(userenum.DeleteUser, user)
	resp.Status(http.StatusOK).
		MessageT(terms.VDeletedSuccessfully, terms.User).
		JSON()
}

// Excel generate excel files based on search
func (p *UserAPI) Excel(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, userenum.Entity)
	var err error

	params.Limit = p.Engine.Environments.ExcelMaxRows
	params.Offset = 0
	params.Order = fmt.Sprintf("%v.id ASC", usermodel.Table)

	users, err := p.Service.Repo.List(params)
	if err != nil {
		err = p.Engine.ErrorLog.TickCustom(err, "E1000147",
			generr.InternalServerErr, "", terms.CanNotGenerateTheExcelList)
		resp.Error(err).JSON()
		return
	}

	ex := excel.New("user").
		AddSheet("users").
		AddSheet("Summary").
		Active("users").
		SetPageLayout("landscape", "A4").
		SetPageMargins(0.2).
		SetHeaderFooter().
		SetColWidth("A", "C", 17).
		SetColWidth("D", "F", 15.3).
		SetColWidth("H", "H", 20).
		SetColWidth("N", "O", 20).
		Active("Summary").
		Active("Nodes").
		WriteHeader("ID", "FullName", "Username", "Role", "Lang", "Email")

	for i, v := range users {
		column := &[]interface{}{
			v.ID,
			v.Name,
			v.Username,
			v.Role,
			v.Lang,
			v.Email,
		}
		err = ex.File.SetSheetRow(ex.ActiveSheet, fmt.Sprint("A", i+2), column)
		p.Engine.ServerLog.CheckError(err, "E1000148", terms.CanNotGenerateTheExcelList, params, users)
	}

	ex.Sheets[ex.ActiveSheet].Row = len(users) + 1

	ex.AddTable()

	buffer, downloadName, err := ex.Generate()
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(userenum.ExcelUser)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+downloadName)
	c.Data(http.StatusOK, "application/octet-stream", buffer.Bytes())

}
