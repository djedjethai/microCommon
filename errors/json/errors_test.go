package json

import (
	"gitlab.com/grpasr/common/tests"
	"testing"
)

// CodeError
func TestCreateNewCodeErrorBasic(t *testing.T) {
	tests.MaybeFail = tests.InitFailFunc(t)

	nec := NewCustomCodeError(ErrInvalidRequest)

	e := nec.Error()

	tests.MaybeFail("create_new_codeError_basic", tests.Expect(e, "1000 : The request is invalid"))
}

func TestCreateNewCodeErrorWithService(t *testing.T) {
	tests.MaybeFail = tests.InitFailFunc(t)

	nec := NewCustomCodeError(ErrUnauthorizedClient, "serviceX")

	e := nec.Error()

	tests.MaybeFail("Create_new_codeError_with_service", tests.Expect(e, "1001 : The client is not authorized to access the requested resource, Service: serviceX"))
}

func TestCreateNewCodeErrorWithServiceAndComment(t *testing.T) {
	tests.MaybeFail = tests.InitFailFunc(t)

	nec := NewCustomCodeError(ErrServerError, "serviceX", "A specific comment")

	e := nec.Error()

	tests.MaybeFail("Create_new_codeError_with_service_comment", tests.Expect(e, "1003 : The server encountered an internal error while processing the request, Service: serviceX, Comment: A specific comment"))
}

func TestCreateNewCodeErrorWithComment(t *testing.T) {
	tests.MaybeFail = tests.InitFailFunc(t)

	nec := NewCustomCodeError(ErrServerError, "", "A specific comment")

	e := nec.Error()

	tests.MaybeFail("Create_new_codeError_with_comment", tests.Expect(e, "1003 : The server encountered an internal error while processing the request, Comment: A specific comment"))
}

func TestCreateNewCodeErrorMarshalJSON(t *testing.T) {
	tests.MaybeFail = tests.InitFailFunc(t)

	nec := NewCustomCodeError(ErrServerError, "", "A specific comment")

	res, err := nec.MarshalJSON()

	tests.MaybeFail("Marshal_json", err, tests.Expect(string(res), `{"code":1003,"description":"The server encountered an internal error while processing the request","comment":"A specific comment"}`))
}

func TestCreateNewCodeErrorUnmarshalJSONBasic(t *testing.T) {
	tests.MaybeFail = tests.InitFailFunc(t)

	nce := NewCodeError(ErrAccessDenied)

	nnn := NewCustomCodeError(ErrInvalidClient)

	res, _ := nnn.MarshalJSON()

	err := nce.UnmarshalJSON(res)

	tests.MaybeFail("Unmarshal_json_basic", err,
		tests.Expect(nce.Code, 1005),
		tests.Expect(nce.Description, "The client is invalid or not recognized"))
}

func TestCreateNewCodeErrorUnmarshalJSONWithService(t *testing.T) {
	tests.MaybeFail = tests.InitFailFunc(t)

	nce := NewCodeError(ErrAccessDenied, "serverAA", "wrongPassword")

	nnn := NewCustomCodeError(ErrInvalidClient, "newService", "newComment")

	res, _ := nnn.MarshalJSON()

	err := nce.UnmarshalJSON(res)

	tests.MaybeFail("Unmarshal_json_with_service", err,
		tests.Expect(nce.Code, 1005),
		tests.Expect(nce.Description, "The client is invalid or not recognized"),
		tests.Expect(nce.Service, "newService"),
		tests.Expect(nce.Comment, "newComment"))
}

func TestCodeErrorReturnCode(t *testing.T) {
	tests.MaybeFail = tests.InitFailFunc(t)

	nce := NewCustomCodeError(ErrAccessDenied, "serverAA", "wrongPassword")

	tests.MaybeFail("Return the error code", tests.Expect(nce.GetCode(), 1002))
}

func TestCreateNewCodeErrorWithPayloadUnmarshalJSONWithService(t *testing.T) {
	tests.MaybeFail = tests.InitFailFunc(t)

	payload := make(map[string]interface{})

	nce := NewCodeError(ErrAccessDenied, "serverAA", "wrongPassword")

	nnn := NewCustomCodeError(ErrInvalidClient, "newService", "newComment")
	payload["string"] = "a string"
	payload["int"] = float64(34)
	nnn.SetPayload(payload)

	res, _ := nnn.MarshalJSON()

	err := nce.UnmarshalJSON(res)

	tests.MaybeFail("Unmarshal_json_with_service", err,
		tests.Expect(nce.Code, 1005),
		tests.Expect(nce.Description, "The client is invalid or not recognized"),
		tests.Expect(nce.Service, "newService"),
		tests.Expect(nce.Comment, "newComment"),
		tests.Expect(nce.GetPayload()["string"], "a string"),
		tests.Expect(nce.GetPayload()["int"].(float64), float64(34)))
}

// statusCode
func TestCreateNewHTTPStatusBasic(t *testing.T) {
	tests.MaybeFail = tests.InitFailFunc(t)

	nec := NewCustomHTTPStatus(StatusOK)

	e := nec.Error()

	tests.MaybeFail("Create_new_HTTPStatus_basic", tests.Expect(e, "200 : Request successful"))
}

func TestCreateHTTPStatusWithURI(t *testing.T) {
	tests.MaybeFail = tests.InitFailFunc(t)

	nec := NewCustomHTTPStatus(StatusNotFound, "/my/path/status")

	e := nec.Error()

	tests.MaybeFail("Create_new_HTTPStatus_with_URI", tests.Expect(e, "404 : Resource not found, Uri: /my/path/status"))
}

func TestCreateHTTPStatusWithURIAndComment(t *testing.T) {
	tests.MaybeFail = tests.InitFailFunc(t)

	nec := NewCustomHTTPStatus(StatusForbidden, "/forbidden/path/ddd", "Private area")

	e := nec.Error()

	tests.MaybeFail("Create_new_HTTPStatus_with_URI_comment", tests.Expect(e, "403 : Request forbidden, Uri: /forbidden/path/ddd, Comment: Private area"))
}

func TestCreateNewHTTPStatusWithComment(t *testing.T) {
	tests.MaybeFail = tests.InitFailFunc(t)

	nec := NewCustomHTTPStatus(StatusInternalServerError, "", "A specific comment")

	e := nec.Error()

	tests.MaybeFail("Create_new_HTTPStatus_comment", tests.Expect(e, "500 : Server error, Comment: A specific comment"))
}

func TestCreateHTTPStatusErrorMarshalJSON(t *testing.T) {
	tests.MaybeFail = tests.InitFailFunc(t)

	nec := NewCustomHTTPStatus(StatusBadRequest, "", "A specific comment")

	res, err := nec.MarshalJSON()

	tests.MaybeFail("Create_new_HTTPStatus_marshalJSON", err, tests.Expect(string(res), `{"code":400,"description":"Invalid request","comment":"A specific comment"}`))
}

func TestCreateNewHTTPStatusUnmarshalJSONBasic(t *testing.T) {
	tests.MaybeFail = tests.InitFailFunc(t)

	nce := NewHTTPStatus(StatusSeeOther)

	nnn := NewCustomHTTPStatus(StatusCreated, "/the/path", "That's ok !!!")

	res, _ := nnn.MarshalJSON()

	err := nce.UnmarshalJSON(res)

	tests.MaybeFail("Unmarshal_json_with_basic", err,
		tests.Expect(nce.Code, 201),
		tests.Expect(nce.Description, "Resource created"),
		tests.Expect(nce.URI, "/the/path"),
		tests.Expect(nce.Comment, "That's ok !!!"),
	)
}

func TestCreateNewHTTPStatusUnmarshalJSONWithService(t *testing.T) {
	tests.MaybeFail = tests.InitFailFunc(t)

	nce := NewHTTPStatus(StatusFound, "/the/path", "Let's try this")

	nnn := NewCustomHTTPStatus(StatusNoContent)

	res, _ := nnn.MarshalJSON()

	err := nce.UnmarshalJSON(res)

	tests.MaybeFail("createNewCodeError", err,
		tests.Expect(nce.Code, 204),
		tests.Expect(nce.Description, "No content to return"),
		tests.Expect(len(nce.URI), 0),
		tests.Expect(len(nce.Comment), 0))
}

func TestStatusCodeReturnCode(t *testing.T) {
	tests.MaybeFail = tests.InitFailFunc(t)

	nec := NewCustomHTTPStatus(StatusOK)

	tests.MaybeFail("getCode", tests.Expect(nec.GetCode(), 200))
}

func TestCreateNewHTTPStatusWithPayloadUnmarshalJSONWithService(t *testing.T) {
	tests.MaybeFail = tests.InitFailFunc(t)

	payload := make(map[string]interface{})

	nce := NewHTTPStatus(StatusFound, "/the/path", "Let's try this")

	nnn := NewCustomHTTPStatus(StatusNoContent)
	payload["string"] = "a string"
	payload["int"] = float64(34)
	nnn.SetPayload(payload)

	res, _ := nnn.MarshalJSON()

	err := nce.UnmarshalJSON(res)

	tests.MaybeFail("createNewCodeError", err,
		tests.Expect(nce.Code, 204),
		tests.Expect(nce.Description, "No content to return"),
		tests.Expect(len(nce.URI), 0),
		tests.Expect(len(nce.Comment), 0),
		tests.Expect(nce.GetPayload()["string"], "a string"),
		tests.Expect(nce.GetPayload()["int"].(float64), float64(34)))
}
