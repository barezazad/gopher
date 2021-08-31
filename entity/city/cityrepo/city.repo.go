package cityrepo

import (
	"gopher/entity/city/citymodel"
	"gopher/internal/core"
	"gopher/internal/core/dberror"
	"gopher/internal/core/terms"
	"gopher/internal/core/validator"
	"gopher/internal/param"
	"gopher/pkg/helper"
	"reflect"

	"gorm.io/gorm/clause"
)

// CityRepo for injecting engine
type CityRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideCityRepo is used in wire and initiate the Cols
func ProvideCityRepo(engine *core.Engine) CityRepo {
	return CityRepo{
		Engine: engine,
		Cols:   helper.TagExtractor(reflect.TypeOf(citymodel.City{}), citymodel.Table),
	}
}

// FindByID finds the city via its id
func (p *CityRepo) FindByID(params param.Param) (city citymodel.City, err error) {

	err = params.GetDB(p.Engine.DB).
		Where("cities.id = ? ", params.ID).First(&city).Error

	err = dberror.DbError(p.Engine, err, "E1000020", city, citymodel.Table, terms.Info)
	return
}

// FindByIDTx finds the city via its id and lock row (for update)
func (p *CityRepo) FindByIDTx(params param.Param) (city citymodel.City, err error) {

	err = params.GetDB(p.Engine.DB).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("cities.id = ? ", params.ID).First(&city).Error

	err = dberror.DbError(p.Engine, err, "E1000021", city, citymodel.Table, terms.Info)
	return
}

// FindAll it return list of all city
func (p *CityRepo) FindAll(params param.Param) (cities []citymodel.City, err error) {

	err = params.GetDB(p.Engine.DB).Find(&cities).Error

	err = dberror.DbError(p.Engine, err, "E1000022", citymodel.City{}, citymodel.Table, terms.Info)
	return
}

// List returns an array of cities
func (p *CityRepo) List(params param.Param) (cities []citymodel.City, err error) {

	var colsStr string
	if colsStr, err = validator.CheckColumns(p.Engine, p.Cols, params.Select, "E1000023"); err != nil {
		return
	}

	var whereStr string
	if whereStr, err = params.ParseWhere(p.Engine, p.Cols, "E1000024"); err != nil {
		return
	}

	err = params.GetDB(p.Engine.DB).
		Select(colsStr).
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&cities).Error

	err = dberror.DbError(p.Engine, err, "E1000025", citymodel.City{}, citymodel.Table, terms.List)

	return
}

// Count of cities, mainly calls with List
func (p *CityRepo) Count(params param.Param) (count int64, err error) {

	var whereStr string
	if whereStr, err = params.ParseWhere(p.Engine, p.Cols, "E1000026"); err != nil {
		return
	}

	err = params.GetDB(p.Engine.DB).Table(citymodel.Table).
		Where(whereStr).
		Count(&count).Error

	err = dberror.DbError(p.Engine, err, "E1000027", citymodel.City{}, citymodel.Table, terms.List)
	return
}

// Create a city
func (p *CityRepo) Create(city citymodel.City, params param.Param) (u citymodel.City, err error) {

	if err = params.GetDB(p.Engine.DB).Create(&city).Scan(&u).Error; err != nil {
		err = dberror.DbError(p.Engine, err, "E1000028", city, citymodel.Table, terms.Created)
	}
	return
}

// Save the city, in case it is not exist create it
func (p *CityRepo) Save(city citymodel.City, params param.Param) (u citymodel.City, err error) {

	if err = params.GetDB(p.Engine.DB).Save(&city).Find(&u).Error; err != nil {
		err = dberror.DbError(p.Engine, err, "E1000029", city, citymodel.Table, terms.Saved)
		return
	}

	return
}

// Delete the city
func (p *CityRepo) Delete(city citymodel.City, params param.Param) (err error) {

	if err = params.GetDB(p.Engine.DB).Unscoped().Delete(&city).Error; err != nil {
		err = dberror.DbError(p.Engine, err, "E1000030", city, citymodel.Table, terms.Deleted)
	}
	return
}
