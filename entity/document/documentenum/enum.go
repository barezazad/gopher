package documentenum

const (
	Entity         = "documents"
	CreateDocument = "document-create"
	UpdateDocument = "document-update"
	DeleteDocument = "document-delete"
	ListDocument   = "document-list"
	ViewDocument   = "document-view"
	AllDocument    = "document-all"
	ExcelDocument  = "document-excel"
)

// DocumentType enums
const (
	Cities = "cities"
	Gifts  = "gifts"
)

// all accepted document types
var DocumentTypes = []string{
	Cities,
	Gifts,
}

// accepted type of documents
const (
	AcceptedImage      string = ".png,.jpeg,.jpg,.svg"
	AcceptedVideo      string = ".mp4,.avi,.mwv,.flv,.mkv"
	AcceptedCoverage   string = ".kml"
	AcceptedVideoImage string = ".png,.jpeg,.jpg,.mp4,.avi,.mwv,.flv,.mkv"
	AcceptedPdfFile    string = ".pdf"
)

// all SingleDocs type
var SingleDocs = []string{
	Gifts,
}

func IsSingleDocs(docType string) (result bool) {
	for _, v := range SingleDocs {
		if v == docType {
			result = true
			return
		}
	}
	return
}
