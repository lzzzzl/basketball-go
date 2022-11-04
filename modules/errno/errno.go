package errno

import "fmt"

var (
	// OK ...
	OK = &Errno{Code: 200, Message: "OK"}

	// Ststem Internal Error
	InternalServerError = &Errno{Code: 10001, Message: "Internal Server Error"}
	ErrBind             = &Errno{Code: 10002, Message: "Input Params Error"}
	ErrTokenSign        = &Errno{Code: 10003, Message: "Signature JWT Error"}
	ErrEncrypt          = &Errno{Code: 10004, Message: "Encrypt Error"}

	// HTTP Error
	ErrHTTPRequestTimeout  = &Errno{Code: 40100, Message: "HTTP Request Timeout"}
	ErrHTTPRequestInternal = &Errno{Code: 40101, Message: "HTTP Request Internal Error"}
)

type Errno struct {
	Code    int
	Message string
}

// Errno Error No
func (e Errno) Error() string {
	return e.Message
}

type Err struct {
	Code    int
	Message string
	Errord  error
}

func (e *Err) Error() string {
	return fmt.Sprintf("Err - code: %d, message: %s, error: %s", e.Code, e.Message, e.Errord)
}

// NewError new error
func NewError(errno *Errno, err error) *Err {
	return &Err{
		Code:    errno.Code,
		Message: errno.Message,
		Errord:  err,
	}
}

func DecoderErr(err error) (int, string) {
	if err == nil {
		return OK.Code, OK.Message
	}
	switch typed := err.(type) {
	case *Err:
		if typed.Code == ErrBind.Code {
			typed.Message = typed.Message + " details is " + typed.Errord.Error()
		}
		return typed.Code, typed.Message
	case *Errno:
		return typed.Code, typed.Message
	default:
	}
	return InternalServerError.Code, err.Error()
}
