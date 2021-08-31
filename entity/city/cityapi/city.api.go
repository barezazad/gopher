package cityapi

import (
	"fmt"
	"gopher/entity/city/cityenum"
	"gopher/entity/city/citymodel"
	"gopher/entity/service"
	"gopher/internal/core"
	"gopher/internal/core/terms"
	"gopher/internal/response"
	"net/http"
	"runtime/debug"

	"gopher/pkg/dictionary"

	"github.com/gin-gonic/gin"
)

// CityAPI for injecting city service
type CityAPI struct {
	Service service.CityService
	Engine  *core.Engine
}

// ProvideCityAPI for city is used in wire
func ProvideCityAPI(c service.CityService) CityAPI {
	return CityAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch a city by it's id
func (p *CityAPI) FindByID(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, citymodel.Table)
	var city citymodel.City
	var err error

	if params.ID, err = resp.GetID(c.Param("cityID"), "E1000015"); err != nil {
		return
	}

	if city, err = p.Service.FindByID(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(cityenum.ViewCity)
	resp.Status(http.StatusOK).
		MessageT(terms.VInfo, terms.City).
		JSON(city)
}

// FindAll it return list of all city
func (p *CityAPI) FindAll(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, citymodel.Table)
	var cities []citymodel.City
	var err error

	if cities, err = p.Service.FindAll(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(cityenum.AllCity)
	resp.Status(http.StatusOK).
		MessageT(terms.AllV, terms.Cities).
		JSON(cities)
}

// List of cities
func (p *CityAPI) List(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, citymodel.Table)
	data := make(map[string]interface{})
	var err error

	if data["list"], data["count"], err = p.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(cityenum.ListCity)
	resp.Status(http.StatusOK).
		MessageT(terms.ListOfV, terms.Cities).
		JSON(data)
}

// Create city
func (p *CityAPI) Create(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, citymodel.Table)
	var city, createdCity citymodel.City
	var err error

	if err = resp.Context.ShouldBind(&city); err != nil {
		resp.NotBind(err, "E1000016")
		return
	}

	// start transaction
	params.Tx = p.Engine.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			p.Engine.ServerLog.Debug("TRANSACTION ROLLBACK", fmt.Errorf(`%v`, r),
				string(debug.Stack()), "create city")
			err = fmt.Errorf(`%v`, r)
			// rollback transaction
			params.Tx.Rollback()
			return
		}
	}()

	city.CreatedBy = params.UserID
	if createdCity, err = p.Service.Create(city, params); err != nil {
		resp.Error(err).JSON()
		// rollback transaction
		params.Tx.Rollback()
		return
	}

	// commit transaction
	params.Tx.Commit()

	resp.RecordCreate(cityenum.CreateCity, createdCity)
	resp.Status(http.StatusOK).
		MessageT(terms.VCreatedSuccessfully, dictionary.Translate(terms.City)).
		JSON(createdCity)
}

// Update city
func (p *CityAPI) Update(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, citymodel.Table)
	var city, cityBefore, cityUpdated citymodel.City
	var err error

	if params.ID, err = resp.GetID(c.Param("cityID"), "E1000017"); err != nil {
		return
	}

	if err = resp.Context.ShouldBind(&city); err != nil {
		resp.NotBind(err, "E1000018")
		return
	}

	// start transaction
	params.Tx = p.Engine.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			p.Engine.ServerLog.Debug("TRANSACTION ROLLBACK", fmt.Errorf(`%v`, r),
				string(debug.Stack()), "update city")
			err = fmt.Errorf(`%v`, r)
			// rollback transaction
			params.Tx.Rollback()
			return
		}
	}()

	city.ID = params.ID
	city.UpdatedBy = params.UserID
	if cityUpdated, cityBefore, err = p.Service.Save(city, params); err != nil {
		resp.Error(err).JSON()
		// rollback transaction
		params.Tx.Rollback()
		return
	}

	// commit transaction
	params.Tx.Commit()

	resp.Record(cityenum.UpdateCity, cityBefore, cityUpdated)
	resp.Status(http.StatusOK).
		MessageT(terms.VUpdatedSuccessfully, dictionary.Translate(terms.City)).
		JSON(cityUpdated)
}

// Delete city
func (p *CityAPI) Delete(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, citymodel.Table)
	var city citymodel.City
	var err error

	if params.ID, err = resp.GetID(c.Param("cityID"), "E1000019"); err != nil {
		return
	}

	// start transaction
	params.Tx = p.Engine.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			p.Engine.ServerLog.Debug("TRANSACTION ROLLBACK", fmt.Errorf(`%v`, r),
				string(debug.Stack()), "delete city")
			err = fmt.Errorf(`%v`, r)
			// rollback transaction
			params.Tx.Rollback()
			return
		}
	}()

	if city, err = p.Service.Delete(params); err != nil {
		resp.Error(err).JSON()
		// rollback transaction
		params.Tx.Rollback()
		return
	}

	// commit transaction
	params.Tx.Commit()

	resp.Record(cityenum.DeleteCity, city)
	resp.Status(http.StatusOK).
		MessageT(terms.VDeletedSuccessfully, terms.City).
		JSON()
}
