package service

import (
	"gopher/entity/document/documentenum"
	"gopher/entity/gift/giftmodel"
	"gopher/entity/gift/giftrepo"
	"gopher/internal/core"
	"gopher/internal/core/validator"
	"gopher/internal/param"
)

type GiftService struct {
	Repo            giftrepo.GiftRepo
	Engine          *core.Engine
	DocumentService DocumentService
}

// ProvideGiftService for gift is used in wire
func ProvideGiftService(p giftrepo.GiftRepo, documentServ DocumentService) GiftService {
	return GiftService{
		Repo:            p,
		Engine:          p.Engine,
		DocumentService: documentServ,
	}
}

// FindByID for getting gift by it's id
func (p *GiftService) FindByID(params param.Param) (gift giftmodel.Gift, err error) {

	if gift, err = p.Repo.FindByID(params); err != nil {
		return
	}

	params.ID = gift.ID
	if gift.Documents, err = p.DocumentService.GetDocsByIdType(params, documentenum.Gifts); err != nil {
		return
	}

	return
}

// List of gifts, it support pagination and search and return back count
func (p *GiftService) List(params param.Param) (gifts []giftmodel.Gift, count int64, err error) {

	if gifts, err = p.Repo.List(params); err != nil {
		return
	}

	for i, v := range gifts {
		params.ID = v.ID
		gifts[i].Documents, _ = p.DocumentService.GetDocsByIdType(params, documentenum.Gifts)
	}

	if count, err = p.Repo.Count(params); err != nil {
		return
	}
	return
}

// Create a gift
func (p *GiftService) Create(gift giftmodel.Gift, params param.Param) (createdGift giftmodel.Gift, err error) {

	// validate to create gift
	if err = validator.BindTagExtractor(p.Engine, gift, "E1000113", giftmodel.Table, core.Create); err != nil {
		return
	}

	if createdGift, err = p.Repo.Create(gift, params); err != nil {
		return
	}

	if gift.Attachments != nil {
		params.ID = createdGift.ID
		if createdGift.Documents, err = p.DocumentService.UploadDocsBase64(params, documentenum.Gifts,
			documentenum.AcceptedImage, gift.Attachments); err != nil {
			return
		}
	}

	return
}

// Save gift
func (p *GiftService) Save(gift giftmodel.Gift, params param.Param) (updatedGift, giftBefore giftmodel.Gift, err error) {

	// validate to update gift
	if err = validator.BindTagExtractor(p.Engine, gift, "E1000114", giftmodel.Table, core.Save); err != nil {
		return
	}

	if giftBefore, err = p.FindByID(params); err != nil {
		return
	}

	if updatedGift, err = p.Repo.Save(gift, params); err != nil {
		return
	}

	if gift.Attachments != nil {
		if updatedGift.Documents, err = p.DocumentService.UploadDocsBase64(params, documentenum.Gifts,
			documentenum.AcceptedImage, gift.Attachments); err != nil {
			return
		}
	}

	return
}

// Delete gift
func (p *GiftService) Delete(params param.Param) (gift giftmodel.Gift, err error) {

	if gift, err = p.FindByID(params); err != nil {
		return
	}

	if err = p.Repo.Delete(gift, params); err != nil {
		return
	}

	if err = p.DocumentService.DeleteAllDocs(params, documentenum.Gifts); err != nil {
		return
	}

	return
}
