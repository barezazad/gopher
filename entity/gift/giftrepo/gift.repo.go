package giftrepo

import (
	"gopher/entity/gift/giftmodel"
	"gopher/internal/core"
	"gopher/internal/core/dberror"
	"gopher/internal/core/terms"
	"gopher/internal/core/validator"
	"gopher/internal/param"
	"gopher/pkg/helper"
	"reflect"

	"gorm.io/gorm/clause"
)

// GiftRepo for injecting engine
type GiftRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideGiftRepo is used in wire and initiate the Cols
func ProvideGiftRepo(engine *core.Engine) GiftRepo {
	return GiftRepo{
		Engine: engine,
		Cols:   helper.TagExtractor(reflect.TypeOf(giftmodel.Gift{}), giftmodel.Table),
	}
}

// FindByID finds the gift via its id
func (p *GiftRepo) FindByID(params param.Param) (gift giftmodel.Gift, err error) {

	err = params.GetDB(p.Engine.DB).
		Where("gifts.id = ? ", params.ID).First(&gift).Error

	err = dberror.DbError(p.Engine, err, "E1000050", gift, giftmodel.Table, terms.Info)
	return
}

// FindByIDTx finds the gift via its id and lock row (for update)
func (p *GiftRepo) FindByIDTx(params param.Param) (gift giftmodel.Gift, err error) {

	err = params.GetDB(p.Engine.DB).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("gifts.id = ? ", params.ID).First(&gift).Error

	err = dberror.DbError(p.Engine, err, "E1000051", gift, giftmodel.Table, terms.Info)
	return
}

// List returns an array of gifts
func (p *GiftRepo) List(params param.Param) (gifts []giftmodel.Gift, err error) {

	var colsStr string
	if colsStr, err = validator.CheckColumns(p.Engine, p.Cols, params.Select, "E1000052"); err != nil {
		return
	}

	var whereStr string
	if whereStr, err = params.ParseWhere(p.Engine, p.Cols, "E1000053"); err != nil {
		return
	}

	err = params.GetDB(p.Engine.DB).
		Select(colsStr).
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&gifts).Error

	err = dberror.DbError(p.Engine, err, "E1000054", giftmodel.Gift{}, giftmodel.Table, terms.List)

	return
}

// Count of gifts, mainly calls with List
func (p *GiftRepo) Count(params param.Param) (count int64, err error) {

	var whereStr string
	if whereStr, err = params.ParseWhere(p.Engine, p.Cols, "E1000055"); err != nil {
		return
	}

	err = params.GetDB(p.Engine.DB).Table(giftmodel.Table).
		Where(whereStr).
		Count(&count).Error

	err = dberror.DbError(p.Engine, err, "E1000056", giftmodel.Gift{}, giftmodel.Table, terms.List)
	return
}

// Create a gift
func (p *GiftRepo) Create(gift giftmodel.Gift, params param.Param) (u giftmodel.Gift, err error) {

	if err = params.GetDB(p.Engine.DB).Create(&gift).Scan(&u).Error; err != nil {
		err = dberror.DbError(p.Engine, err, "E1000057", gift, giftmodel.Table, terms.Created)
	}
	return
}

// Save the gift, in case it is not exist create it
func (p *GiftRepo) Save(gift giftmodel.Gift, params param.Param) (u giftmodel.Gift, err error) {

	if err = params.GetDB(p.Engine.DB).Save(&gift).Find(&u).Error; err != nil {
		err = dberror.DbError(p.Engine, err, "E1000058", gift, giftmodel.Table, terms.Saved)
		return
	}

	return
}

// Delete the gift
func (p *GiftRepo) Delete(gift giftmodel.Gift, params param.Param) (err error) {

	if err = params.GetDB(p.Engine.DB).Unscoped().Delete(&gift).Error; err != nil {
		err = dberror.DbError(p.Engine, err, "E1000059", gift, giftmodel.Table, terms.Deleted)
	}
	return
}
