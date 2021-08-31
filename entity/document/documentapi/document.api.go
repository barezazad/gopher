package documentapi

import (
	"gopher/entity/document/documentenum"
	"gopher/entity/document/documentmodel"
	"gopher/entity/service"
	"gopher/internal/core"
	"gopher/internal/core/terms"
	"gopher/internal/response"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// DocumentAPI for injecting document service
type DocumentAPI struct {
	Service service.DocumentService
	Engine  *core.Engine
}

// ProvideDocumentAPI for document is used in wire
func ProvideDocumentAPI(c service.DocumentService) DocumentAPI {
	return DocumentAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch a document by it's id
func (p *DocumentAPI) FindByID(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, documentmodel.Table)
	var document documentmodel.Document
	var err error

	if params.ID, err = resp.GetID(c.Param("documentID"), "E1000031"); err != nil {
		return
	}

	if document, err = p.Service.FindByID(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(documentenum.ViewDocument)
	resp.Status(http.StatusOK).
		MessageT(terms.VInfo, terms.Document).
		JSON(document)
}

// List of documents
func (p *DocumentAPI) List(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, documentmodel.Table)
	data := make(map[string]interface{})
	var err error

	if data["list"], data["count"], err = p.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(documentenum.ListDocument)
	resp.Status(http.StatusOK).
		MessageT(terms.ListOfV, terms.Documents).
		JSON(data)
}

// DownloadDocs finds the document via its doc_id and doc_type
func (p *DocumentAPI) DownloadDocs(c *gin.Context) {
	var err error
	docName := c.Param("docName")
	docType := c.Param("docType")
	var folderName string

	if folderName, err = p.Service.DocumentLocationFolder(docType); err != nil {
		return
	}

	fileFullPath := filepath.Join(folderName, docName)
	if _, err = os.Stat(fileFullPath); os.IsNotExist(err) {
		p.Engine.ServerLog.CheckError(err, "can't download document", fileFullPath)
		// default image
		//fileFullPath = filepath.Join("assets/img/", docName)
	}
	c.FileAttachment(fileFullPath, docName)

}

// Delete document
func (p *DocumentAPI) Delete(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, documentmodel.Table)
	var document documentmodel.Document
	var err error

	docType := c.Param("docType")
	if params.ID, err = resp.GetID(c.Param("documentID"), "E1000032"); err != nil {
		return
	}

	var folderName string
	folderName, _ = p.Service.DocumentLocationFolder(docType)

	// delete in database
	if document, err = p.Service.Delete(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	// delete in system
	oldFile := filepath.Join(folderName, document.Name)
	os.Remove(oldFile)

	resp.Record(documentenum.DeleteDocument, document)
	resp.Status(http.StatusOK).
		MessageT(terms.VDeletedSuccessfully, terms.Document).
		JSON()
}
