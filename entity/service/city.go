package service

import (
	"encoding/json"
	"gopher/entity/city/citymodel"
	"gopher/entity/city/cityrepo"
	"gopher/entity/document/documentenum"
	"gopher/internal/core"
	"gopher/internal/core/validator"
	"gopher/internal/param"
)

type CityService struct {
	Repo            cityrepo.CityRepo
	Engine          *core.Engine
	DocumentService DocumentService
}

// ProvideCityService for city is used in wire
func ProvideCityService(p cityrepo.CityRepo, documentServ DocumentService) CityService {
	return CityService{
		Repo:            p,
		Engine:          p.Engine,
		DocumentService: documentServ,
	}
}

// FindByID for getting city by it's id
func (p *CityService) FindByID(params param.Param) (city citymodel.City, err error) {

	if city, err = p.Repo.FindByID(params); err != nil {
		return
	}

	params.ID = city.ID
	if city.Documents, err = p.DocumentService.GetDocsByIdType(params, documentenum.Cities); err != nil {
		return
	}

	return
}

// FindAll it return list of all city
func (p *CityService) FindAll(params param.Param) (cities []citymodel.City, err error) {

	if cacheCities, ok := p.getInCache(); ok {
		cities = cacheCities
		return
	}

	if cities, err = p.Repo.FindAll(params); err != nil {
		return
	}

	for i, v := range cities {
		params.ID = v.ID
		cities[i].Documents, _ = p.DocumentService.GetDocsByIdType(params, documentenum.Cities)
	}

	// save all cities in cache
	p.Engine.Cache.Set(citymodel.Table, cities, p.Engine.Environments.Redis.CacheApiTTL)
	return
}

// List of cities, it support pagination and search and return back count
func (p *CityService) List(params param.Param) (cities []citymodel.City, count int64, err error) {

	if cities, err = p.Repo.List(params); err != nil {
		return
	}

	for i, v := range cities {
		params.ID = v.ID
		cities[i].Documents, _ = p.DocumentService.GetDocsByIdType(params, documentenum.Cities)
	}

	if count, err = p.Repo.Count(params); err != nil {
		return
	}
	return
}

// Create a city
func (p *CityService) Create(city citymodel.City, params param.Param) (createdCity citymodel.City, err error) {

	// validate to create city
	if err = validator.BindTagExtractor(p.Engine, city, "E1000097", citymodel.Table, core.Create); err != nil {
		return
	}

	if createdCity, err = p.Repo.Create(city, params); err != nil {
		return
	}

	if city.Attachments != nil {
		params.ID = createdCity.ID
		if createdCity.Documents, err = p.DocumentService.UploadDocs(params, documentenum.Cities,
			documentenum.AcceptedImage, city.Attachments); err != nil {
			return
		}
	}

	// remove cities in cache
	p.Engine.Cache.Delete(citymodel.Table)

	return
}

// Save city
func (p *CityService) Save(city citymodel.City, params param.Param) (updatedCity, cityBefore citymodel.City, err error) {

	// validate to update city
	if err = validator.BindTagExtractor(p.Engine, city, "E1000098", citymodel.Table, core.Save); err != nil {
		return
	}

	if cityBefore, err = p.FindByID(params); err != nil {
		return
	}

	if updatedCity, err = p.Repo.Save(city, params); err != nil {
		return
	}

	if city.Attachments != nil {
		if updatedCity.Documents, err = p.DocumentService.UploadDocs(params, documentenum.Cities,
			documentenum.AcceptedImage, city.Attachments); err != nil {
			return
		}
	}

	// remove cities in cache
	p.Engine.Cache.Delete(citymodel.Table)

	return
}

// Delete city
func (p *CityService) Delete(params param.Param) (city citymodel.City, err error) {

	if city, err = p.FindByID(params); err != nil {
		return
	}

	if err = p.Repo.Delete(city, params); err != nil {
		return
	}

	if err = p.DocumentService.DeleteAllDocs(params, documentenum.Cities); err != nil {
		return
	}

	// remove cities in cache
	p.Engine.Cache.Delete(citymodel.Table)

	return
}

// check value in cache and get it , if exist
func (p *CityService) getInCache() (cities []citymodel.City, ok bool) {

	result, err := p.Engine.Cache.Get(citymodel.Table)
	if err != nil {
		return
	}

	if result != "" {
		if err = json.Unmarshal([]byte(result), &cities); err != nil {
			return
		}
		ok = true
	}

	return
}
