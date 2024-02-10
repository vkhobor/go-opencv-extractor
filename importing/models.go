package importing

type Status string

const (
	StatusZero      Status = "nil"
	StatusImporting        = "importing"
	StatusImported         = "imported"
	StatusError            = "error"
)

type DbEntry struct {
	Status      Status
	Title       string
	Url         string
	FileNames   []string
	ErrorString string
}
