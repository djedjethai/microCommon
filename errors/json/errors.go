package json

import (
	"encoding/json"
	"fmt"
)

type IError interface {
	Error() string
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
	GetCode() int
	SetPayload(map[string]interface{})
	GetPayload() map[string]interface{}
}

// CustomeError is the only Error type use in the code
type CustomError struct {
	IError
}

func NewCustomCodeError(code ErrorCode, infos ...string) CustomError {
	return CustomError{NewCodeError(code, infos...)}
}

func NewCustomHTTPStatus(code StatusCode, uri ...string) CustomError {
	return CustomError{NewHTTPStatus(code, uri...)}
}

func (ce CustomError) Error() string {
	return ce.IError.Error()
}

func (ce CustomError) MarshalJSON() ([]byte, error) {
	return ce.IError.MarshalJSON()
}

func (ce CustomError) UnmarshalJSON(b []byte) error {
	return ce.IError.UnmarshalJSON(b)
}

func (ce CustomError) GetCode() int {
	return ce.IError.GetCode()
}

func (ce CustomError) SetPayload(pl map[string]interface{}) {
	ce.IError.SetPayload(pl)
}

func (ce CustomError) GetPayload() map[string]interface{} {
	return ce.IError.GetPayload()
}

type BaseError struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

type CodeError struct {
	BaseError
	Service string                 `json:"service,omitempty"`
	Comment string                 `json:"comment,omitempty"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}

func NewCodeError(code ErrorCode, srv ...string) *CodeError {
	ce := &CodeError{
		BaseError: BaseError{
			Code:        int(code),
			Description: ErrorCodeDescriptions[code],
		},
	}
	if len(srv) == 1 {
		ce.Service = srv[0]
	}
	if len(srv) > 1 {
		ce.Service = srv[0]
		ce.Comment = srv[1]
	}
	return ce
}

// Error implements the error.Error interface
func (c *CodeError) Error() string {
	if len(c.Service) > 0 && len(c.Comment) == 0 {
		return fmt.Sprintf("%v : %v, Service: %s", c.Code, c.Description, c.Service)
	} else if len(c.Service) > 0 && len(c.Comment) > 0 {
		return fmt.Sprintf("%v : %v, Service: %s, Comment: %s", c.Code, c.Description, c.Service, c.Comment)
	} else if len(c.Service) == 0 && len(c.Comment) > 0 {
		return fmt.Sprintf("%v : %v, Comment: %s", c.Code, c.Description, c.Comment)
	} else {
		return fmt.Sprintf("%v : %s", c.Code, c.Description)
	}
}

// MarshalJSON implements the json.Marshaler interface
func (c *CodeError) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		BaseError
		Service string                 `json:"service,omitempty"`
		Comment string                 `json:"comment,omitempty"`
		Payload map[string]interface{} `json:"payload,omitempty"`
	}{

		c.BaseError,
		c.Service,
		c.Comment,
		c.Payload,
	})
}

// UnmarshalJSON implements the json.Unmarshaller interface
func (c *CodeError) UnmarshalJSON(payload []byte) error {
	var tmp struct {
		BaseError
		Service string                 `json:"service,omitempty"`
		Comment string                 `json:"comment,omitempty"`
		Payload map[string]interface{} `json:"payload,omitempty"`
	}

	var err error
	err = json.Unmarshal(payload, &tmp)

	c.BaseError = tmp.BaseError
	c.Service = tmp.Service
	c.Comment = tmp.Comment
	c.Payload = tmp.Payload

	return err
}

func (c *CodeError) GetCode() int {
	return c.Code
}

func (c *CodeError) SetPayload(pl map[string]interface{}) {
	if pl != nil && len(pl) > 0 {
		c.Payload = pl
	}
}

func (c *CodeError) GetPayload() map[string]interface{} {
	return c.Payload
}

/******
* HTTP STATUS
******/
type HTTPStatus struct {
	BaseError
	URI     string                 `json:"uri,omitempty"`
	Comment string                 `json:"comment,omitempty"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}

func NewHTTPStatus(code StatusCode, infos ...string) *HTTPStatus {
	he := &HTTPStatus{
		BaseError: BaseError{
			Code:        int(code),
			Description: HTTPCodeDescriptions[code],
		},
	}
	if len(infos) > 0 {
		he.URI = infos[0]
	}
	if len(infos) > 1 {
		he.URI = infos[0]
		he.Comment = infos[1]
	}

	return he
}

// Error implements the error.Error interface
func (h *HTTPStatus) Error() string {
	if len(h.URI) > 0 && len(h.Comment) == 0 {
		return fmt.Sprintf("%v : %v, Uri: %v", h.Code, h.Description, h.URI)
	} else if len(h.URI) == 0 && len(h.Comment) > 0 {
		return fmt.Sprintf("%v : %v, Comment: %v", h.Code, h.Description, h.Comment)
	} else if len(h.URI) > 0 && len(h.Comment) > 0 {
		return fmt.Sprintf("%v : %v, Uri: %v, Comment: %v", h.Code, h.Description, h.URI, h.Comment)
	} else {
		return fmt.Sprintf("%v : %s", h.Code, h.Description)
	}
}

// MarshalJSON implements the json.Marshaler interface
func (h *HTTPStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		BaseError
		URI     string                 `json:"uri,omitempty"`
		Comment string                 `json:"comment,omitempty"`
		Payload map[string]interface{} `json:"payload,omitempty"`
	}{
		h.BaseError,
		h.URI,
		h.Comment,
		h.Payload,
	})
}

// UnmarshalJSON implements the json.Unmarshaller interface
func (h *HTTPStatus) UnmarshalJSON(payload []byte) error {
	var err error
	var tmp struct {
		BaseError
		URI     string                 `json:"uri,omitempty"`
		Comment string                 `json:"comment,omitempty"`
		Payload map[string]interface{} `json:"payload,omitempty"`
	}

	err = json.Unmarshal(payload, &tmp)

	h.BaseError = tmp.BaseError
	h.URI = tmp.URI
	h.Comment = tmp.Comment
	h.Payload = tmp.Payload

	return err
}

func (h *HTTPStatus) GetCode() int {
	return h.Code
}

func (h *HTTPStatus) SetPayload(pl map[string]interface{}) {
	if pl != nil && len(pl) > 0 {
		h.Payload = pl
	}
}

func (h *HTTPStatus) GetPayload() map[string]interface{} {
	return h.Payload
}
