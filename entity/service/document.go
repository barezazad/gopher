package service

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"gopher/entity/document/documentenum"
	"gopher/entity/document/documentmodel"
	"gopher/entity/document/documentrepo"
	"gopher/internal/core"
	"gopher/internal/core/terms"
	"gopher/internal/core/validator"
	"gopher/internal/param"
	"gopher/pkg/generr"
	"gopher/pkg/helper/random"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type DocumentService struct {
	Repo   documentrepo.DocumentRepo
	Engine *core.Engine
}

// ProvideDocumentService for document is used in wire
func ProvideDocumentService(p documentrepo.DocumentRepo) DocumentService {
	return DocumentService{
		Repo:   p,
		Engine: p.Engine,
	}
}

// FindByID for getting document by it's id
func (p *DocumentService) FindByID(params param.Param) (document documentmodel.Document, err error) {

	if document, err = p.Repo.FindByID(params); err != nil {
		return
	}
	return
}

// GetDocsByIdType finds the document via its doc_id and doc_type
func (p *DocumentService) GetDocsByIdType(params param.Param, docType string) (document []documentmodel.Document, err error) {

	if document, err = p.Repo.GetDocsByIdType(params, docType); err != nil {
		return
	}
	return
}

// List of documents, it support pagination and search and return back count
func (p *DocumentService) List(params param.Param) (documents []documentmodel.Document, count int64, err error) {

	if documents, err = p.Repo.List(params); err != nil {
		return
	}

	if count, err = p.Repo.Count(params); err != nil {
		return
	}
	return
}

// Create a document
func (p *DocumentService) Create(document documentmodel.Document, params param.Param) (createdDocument documentmodel.Document, err error) {

	// validate to create document
	if err = validator.BindTagExtractor(p.Engine, document, "E1000099", documentmodel.Table, core.Create); err != nil {
		return
	}

	if createdDocument, err = p.Repo.Create(document, params); err != nil {
		return
	}

	return
}

// Save document
func (p *DocumentService) Save(document documentmodel.Document, params param.Param) (updatedDocument, documentBefore documentmodel.Document, err error) {

	// validate to update document
	if err = validator.BindTagExtractor(p.Engine, document, "E1000100", documentmodel.Table, core.Save); err != nil {
		return
	}

	if documentBefore, err = p.FindByID(params); err != nil {
		return
	}

	if updatedDocument, err = p.Repo.Save(document, params); err != nil {
		return
	}

	return
}

// Delete document
func (p *DocumentService) Delete(params param.Param) (document documentmodel.Document, err error) {

	if document, err = p.FindByID(params); err != nil {
		return
	}

	if err = p.Repo.Delete(document, params); err != nil {
		return
	}

	return
}

// UploadDocs to upload documents
func (p *DocumentService) UploadDocs(params param.Param, docType string, acceptedType string,
	docs []*multipart.FileHeader) (documents []documentmodel.Document, err error) {

	var folderName string
	var docPayload documentmodel.Document

	if folderName, err = p.DocumentLocationFolder(docType); err != nil {
		return
	}

	for i, v := range docs {

		// check to validate extention
		fileExt := filepath.Ext(v.Filename)
		if !strings.Contains(acceptedType, fileExt) {
			err = p.Engine.ErrorLog.TickCustom(err, "E1000101", generr.BadRequestErr,
				"", terms.DocumentTypeNotAccepted)
			return
		}
		FileName := random.String(10)
		FileName = fmt.Sprintf(`%v-%v-%v%v`, params.ID, i, FileName, fileExt)
		FilePath := fmt.Sprintf(`%v/%v`, folderName, FileName)

		// read file like io file
		out, _ := v.Open()

		// read file in buffer
		buff := bytes.Buffer{}
		if _, err = buff.ReadFrom(out); err != nil {
			p.Engine.ErrorLog.CheckError(err, "E1000102", "can't read file", buff)
		}

		// check to know this document is single or its allow to multiple
		isSingleDocument := documentenum.IsSingleDocs(docType)
		if isSingleDocument {
			count, _ := p.Repo.CountByIdType(params, docType)
			if count != 0 {
				err = p.Engine.ErrorLog.TickCustom(err, "E1000103", generr.BadRequestErr,
					"", terms.CanNotUploadMultipleDocument)
				return
			}
		}

		// write file
		if err = ioutil.WriteFile(FilePath, buff.Bytes(), 0644); err != nil {
			err = p.Engine.ErrorLog.TickCustom(err, "E1000104", generr.BadRequestErr,
				"", terms.CouldNotUploadDocument)
			return
		}

		docPayload.DocId = params.ID
		docPayload.DocType = docType
		docPayload.Name = FileName
		docPayload.FileType = fileExt

		_, err = p.Create(docPayload, params)
		if err != nil {
			err = p.Engine.ErrorLog.TickCustom(err, "E1000105", generr.BadRequestErr,
				"", terms.DocumentDidNotSaved)
			return
		}
	}

	if documents, err = p.GetDocsByIdType(params, docType); err != nil {
		return
	}

	return
}

// UploadDocsBase64 to upload documents in base64
func (p *DocumentService) UploadDocsBase64(params param.Param, docType string, acceptedType string,
	docs []string) (documents []documentmodel.Document, err error) {

	var folderName string
	var docPayload documentmodel.Document

	if folderName, err = p.DocumentLocationFolder(docType); err != nil {
		return
	}

	for i, v := range docs {

		// decode base64 file
		idx := strings.Index(v, ";base64,")
		reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(v[idx+8:]))

		// read file and append in buffer
		buff := bytes.Buffer{}
		if _, err = buff.ReadFrom(reader); err != nil {
			p.Engine.ErrorLog.CheckError(err, "E1000106", "can't read file", buff)
		}

		// get file extention
		fileExt := fmt.Sprintf(".%v", strings.Split(strings.Split(v, "/")[1], ";")[0])

		// generate file name and path of file
		FileName := random.String(10)
		FileName = fmt.Sprintf(`%v-%v-%v%v`, params.ID, i, FileName, fileExt)
		FilePath := fmt.Sprintf(`%v/%v`, folderName, FileName)

		// check to validate extention
		if !strings.Contains(acceptedType, fileExt) {
			err = p.Engine.ErrorLog.TickCustom(err, "E1000107", generr.BadRequestErr,
				"", terms.DocumentTypeNotAccepted)
			return
		}

		// check to know this document is single or its allow to multiple
		isSingleDocument := documentenum.IsSingleDocs(docType)
		if isSingleDocument {
			count, _ := p.Repo.CountByIdType(params, docType)
			if count != 0 {
				err = p.Engine.ErrorLog.TickCustom(err, "E1000108", generr.BadRequestErr,
					"", terms.CanNotUploadMultipleDocument)
				return
			}
		}

		// write file
		if err = ioutil.WriteFile(FilePath, buff.Bytes(), 0644); err != nil {
			err = p.Engine.ErrorLog.TickCustom(err, "E1000109", generr.BadRequestErr,
				"", terms.CouldNotUploadDocument)
			return
		}

		docPayload.DocId = params.ID
		docPayload.DocType = docType
		docPayload.Name = FileName
		docPayload.FileType = fileExt

		_, err = p.Create(docPayload, params)
		if err != nil {
			err = p.Engine.ErrorLog.TickCustom(err, "E1000110", generr.BadRequestErr,
				"", terms.DocumentDidNotSaved)
			return
		}
	}

	if documents, err = p.GetDocsByIdType(params, docType); err != nil {
		return
	}

	return
}

// DeleteAllDocsById to delete documents
func (p *DocumentService) DeleteAllDocs(params param.Param, docType string) (err error) {

	var folderName string
	var documents []documentmodel.Document

	if folderName, err = p.DocumentLocationFolder(docType); err != nil {
		return
	}

	if documents, err = p.GetDocsByIdType(params, docType); err != nil {
		return
	}

	// delete in database
	for _, v := range documents {
		if err = p.Repo.Delete(v, params); err != nil {
			err = p.Engine.ErrorLog.TickCustom(err, "E1000111", generr.BadRequestErr,
				"", terms.CouldNotDeleteDocuments)
			return
		}
	}

	// delete in system
	for _, v := range documents {
		oldFile := filepath.Join(folderName, v.Name)
		os.Remove(oldFile)
	}

	return
}

// DocumentLocationFolder find folder location by checking type of doc
func (p *DocumentService) DocumentLocationFolder(docType string) (location string, err error) {

	switch docType {
	case documentenum.Cities:
		location = p.Engine.Environments.Document.CitiesDir
		return
	case documentenum.Gifts:
		location = p.Engine.Environments.Document.GiftsDir
		return
	default:
		err = p.Engine.ErrorLog.TickCustom(err, "E1000112", generr.NotFoundErr,
			"", terms.DocumentTypeIsWrong)
		p.Engine.ServerLog.CheckError(err, "document type is wrong,cant find location", docType)
		return
	}
}
