package jsonschema

type Format string

const (
	FormatUnknown  Format = ""
	FormatDateTime Format = "date-time"
	FormatUUID     Format = "uuid"
	FormatEmail    Format = "email"
	FormatURI      Format = "uri"
	FormatDate     Format = "date"
)
