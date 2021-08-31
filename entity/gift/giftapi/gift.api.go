package giftapi

import (
	"fmt"
	"gopher/entity/gift/giftenum"
	"gopher/entity/gift/giftmodel"
	"gopher/entity/service"
	"gopher/internal/core"
	"gopher/internal/core/terms"
	"gopher/internal/response"
	"net/http"
	"runtime/debug"

	"gopher/pkg/dictionary"

	"github.com/gin-gonic/gin"
)

// GiftAPI for injecting gift service
type GiftAPI struct {
	Service service.GiftService
	Engine  *core.Engine
}

// ProvideGiftAPI for gift is used in wire
func ProvideGiftAPI(c service.GiftService) GiftAPI {
	return GiftAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch a gift by it's id
func (p *GiftAPI) FindByID(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, giftmodel.Table)
	var gift giftmodel.Gift
	var err error

	if params.ID, err = resp.GetID(c.Param("giftID"), "E1000045"); err != nil {
		return
	}

	if gift, err = p.Service.FindByID(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(giftenum.ViewGift)
	resp.Status(http.StatusOK).
		MessageT(terms.VInfo, terms.Gift).
		JSON(gift)
}

// List of gifts
func (p *GiftAPI) List(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, giftmodel.Table)
	data := make(map[string]interface{})
	var err error

	if data["list"], data["count"], err = p.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(giftenum.ListGift)
	resp.Status(http.StatusOK).
		MessageT(terms.ListOfV, terms.Gifts).
		JSON(data)
}

// Create gift
func (p *GiftAPI) Create(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, giftmodel.Table)
	var gift, createdGift giftmodel.Gift
	var err error

	if err = resp.BindCipher(&gift, "E1000046"); err != nil {
		return
	}

	// start transaction
	params.Tx = p.Engine.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			p.Engine.ServerLog.Debug("TRANSACTION ROLLBACK", fmt.Errorf(`%v`, r),
				string(debug.Stack()), "create gift")
			err = fmt.Errorf(`%v`, r)
			// rollback transaction
			params.Tx.Rollback()
			return
		}
	}()

	gift.CreatedBy = params.UserID
	if createdGift, err = p.Service.Create(gift, params); err != nil {
		resp.Error(err).JSON()
		// rollback transaction
		params.Tx.Rollback()
		return
	}

	// commit transaction
	params.Tx.Commit()

	resp.RecordCreate(giftenum.CreateGift, createdGift)
	resp.Status(http.StatusOK).
		MessageT(terms.VCreatedSuccessfully, dictionary.Translate(terms.Gift)).
		JSON(createdGift)
}

// Update gift
func (p *GiftAPI) Update(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, giftmodel.Table)
	var gift, giftBefore, giftUpdated giftmodel.Gift
	var err error

	if params.ID, err = resp.GetID(c.Param("giftID"), "E1000047"); err != nil {
		return
	}

	if err = resp.BindCipher(&gift, "E1000048"); err != nil {
		return
	}

	// start transaction
	params.Tx = p.Engine.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			p.Engine.ServerLog.Debug("TRANSACTION ROLLBACK", fmt.Errorf(`%v`, r),
				string(debug.Stack()), "update gift")
			err = fmt.Errorf(`%v`, r)
			// rollback transaction
			params.Tx.Rollback()
			return
		}
	}()

	gift.ID = params.ID
	gift.UpdatedBy = params.UserID
	if giftUpdated, giftBefore, err = p.Service.Save(gift, params); err != nil {
		resp.Error(err).JSON()
		// rollback transaction
		params.Tx.Rollback()
		return
	}

	// commit transaction
	params.Tx.Commit()

	resp.Record(giftenum.UpdateGift, giftBefore, giftUpdated)
	resp.Status(http.StatusOK).
		MessageT(terms.VUpdatedSuccessfully, dictionary.Translate(terms.Gift)).
		JSON(giftUpdated)
}

// Delete gift
func (p *GiftAPI) Delete(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, giftmodel.Table)
	var gift giftmodel.Gift
	var err error

	if params.ID, err = resp.GetID(c.Param("giftID"), "E1000049"); err != nil {
		return
	}

	// start transaction
	params.Tx = p.Engine.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			p.Engine.ServerLog.Debug("TRANSACTION ROLLBACK", fmt.Errorf(`%v`, r),
				string(debug.Stack()), "delete gift")
			err = fmt.Errorf(`%v`, r)
			// rollback transaction
			params.Tx.Rollback()
			return
		}
	}()

	if gift, err = p.Service.Delete(params); err != nil {
		resp.Error(err).JSON()
		// rollback transaction
		params.Tx.Rollback()
		return
	}

	// commit transaction
	params.Tx.Commit()

	resp.Record(giftenum.DeleteGift, gift)
	resp.Status(http.StatusOK).
		MessageT(terms.VDeletedSuccessfully, terms.Gift).
		JSON()
}
