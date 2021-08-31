package documentrepo

import (
	"gopher/entity/document/documentmodel"
	"gopher/internal/core"
	"gopher/internal/core/dberror"
	"gopher/internal/core/terms"
	"gopher/internal/core/validator"
	"gopher/internal/param"
	"gopher/pkg/helper"
	"reflect"

	"gorm.io/gorm/clause"
)

// DocumentRepo for injecting engine
type DocumentRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideDocumentRepo is used in wire and initiate the Cols
func ProvideDocumentRepo(engine *core.Engine) DocumentRepo {
	return DocumentRepo{
		Engine: engine,
		Cols:   helper.TagExtractor(reflect.TypeOf(documentmodel.Document{}), documentmodel.Table),
	}
}

// FindByID finds the document via its id
func (p *DocumentRepo) FindByID(params param.Param) (document documentmodel.Document, err error) {

	err = params.GetDB(p.Engine.DB).
		Where("documents.id = ? ", params.ID).First(&document).Error

	err = dberror.DbError(p.Engine, err, "E1000033", document, documentmodel.Table, terms.Info)
	return
}

// FindByIDTx finds the document via its id
func (p *DocumentRepo) FindByIDTx(params param.Param) (document documentmodel.Document, err error) {

	err = params.GetDB(p.Engine.DB).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("documents.id = ? ", params.ID).First(&document).Error

	err = dberror.DbError(p.Engine, err, "E1000034", document, documentmodel.Table, terms.Info)
	return
}

// CountByIdType count the documents by doc type
func (p *DocumentRepo) CountByIdType(params param.Param, docType string) (count int64, err error) {

	err = params.GetDB(p.Engine.DB).Table(documentmodel.Table).
		Where("documents.doc_id = ? AND documents.doc_type = ? ", params.ID, docType).Count(&count).Error

	err = dberror.DbError(p.Engine, err, "E1000035", documentmodel.Document{}, documentmodel.Table, terms.Info)
	return
}

// GetDocsByIdType finds the document via its doc_id and doc_type
func (p *DocumentRepo) GetDocsByIdType(params param.Param, docType string) (documents []documentmodel.Document, err error) {

	err = params.GetDB(p.Engine.DB).
		Where("documents.doc_id = ? AND documents.doc_type = ? ", params.ID, docType).
		Order("id desc").Find(&documents).Error

	err = dberror.DbError(p.Engine, err, "E1000036", documentmodel.Document{}, documentmodel.Table, terms.Info)
	return
}

// List returns an array of documents
func (p *DocumentRepo) List(params param.Param) (documents []documentmodel.Document, err error) {

	var colsStr string
	if colsStr, err = validator.CheckColumns(p.Engine, p.Cols, params.Select, "E1000037"); err != nil {
		return
	}

	var whereStr string
	if whereStr, err = params.ParseWhere(p.Engine, p.Cols, "E1000038"); err != nil {
		return
	}

	err = params.GetDB(p.Engine.DB).
		Select(colsStr).
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&documents).Error

	err = dberror.DbError(p.Engine, err, "E1000039", documentmodel.Document{}, documentmodel.Table, terms.List)

	return
}

// Count of documents, mainly calls with List
func (p *DocumentRepo) Count(params param.Param) (count int64, err error) {

	var whereStr string
	if whereStr, err = params.ParseWhere(p.Engine, p.Cols, "E1000040"); err != nil {
		return
	}

	err = params.GetDB(p.Engine.DB).Table(documentmodel.Table).
		Where(whereStr).
		Count(&count).Error

	err = dberror.DbError(p.Engine, err, "E1000041", documentmodel.Document{}, documentmodel.Table, terms.List)
	return
}

// Create a document
func (p *DocumentRepo) Create(document documentmodel.Document, params param.Param) (u documentmodel.Document, err error) {

	if err = params.GetDB(p.Engine.DB).Create(&document).Scan(&u).Error; err != nil {
		err = dberror.DbError(p.Engine, err, "E1000042", document, documentmodel.Table, terms.Created)
	}
	return
}

// Save the document, in case it is not exist create it
func (p *DocumentRepo) Save(document documentmodel.Document, params param.Param) (u documentmodel.Document, err error) {

	if err = params.GetDB(p.Engine.DB).Save(&document).Find(&u).Error; err != nil {
		err = dberror.DbError(p.Engine, err, "E1000043", document, documentmodel.Table, terms.Saved)
		return
	}

	return
}

// Delete the document
func (p *DocumentRepo) Delete(document documentmodel.Document, params param.Param) (err error) {

	if err = params.GetDB(p.Engine.DB).Unscoped().Delete(&document).Error; err != nil {
		err = dberror.DbError(p.Engine, err, "E1000044", document, documentmodel.Table, terms.Deleted)
	}
	return
}

// Delete the document by doc type
func (p *DocumentRepo) DeleteByDocType(params param.Param, docType string) (err error) {

	params.GetDB(p.Engine.DB).Exec("DELETE FROM documents where doc_type = ?", docType)
	return
}
