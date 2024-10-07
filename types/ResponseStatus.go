package types

type ResponseStatus int

const (
	Success ResponseStatus = 0
	ParseError ResponseStatus = 1
	InvalidLevelError ResponseStatus = 2
	InvalidSaveError ResponseStatus = 3
	TooLongError ResponseStatus = 4
	InternalError ResponseStatus = 5
)

func (r ResponseStatus) String() string {
	switch r {
	case Success:
		return "SUCCESS"
	case ParseError:
		return "PARSE_ERROR"
	case InvalidLevelError:
		return "INVALID_LEVEL_ERROR"
	case InvalidSaveError:
		return "INVALID_SAVE_ERROR"
	case TooLongError:
		return "TOO_LONG_ERROR"
	case InternalError:
		return "INTERNAL_ERROR"
	default:
		return "UNKNOWN"
	}
}