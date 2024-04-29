package restclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	e "gitlab.com/grpasr/common/errors/json"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type RequestType string

const (
	THandleMultipartWriter RequestType = "THandleMultipartWriter"
	THandleRequest         RequestType = "THandleRequest"
)

// HandleRetryRequest is a generic func that add retry logic to the requests' handlers
func (rs *restService) HandleRetryRequest(ctx context.Context, request *Api, response interface{}, retries int8, delay int8, requestType RequestType, arguments ...string) e.IError {

	baseDelay, _ := time.ParseDuration(fmt.Sprintf("%vs", delay))

	var err e.IError
	for r := int8(0); ; r++ {
		switch requestType {
		case THandleMultipartWriter:
			if len(arguments) > 0 {
				err = rs.HandleMultipartWriter(request, arguments[0], response)
				if err == nil || r >= retries {
					return err
				}
			} else {
				return e.NewCustomHTTPStatus(e.StatusBadRequest, "", "path to write is missing")
			}
		case THandleRequest:
			err = rs.HandleRequest(request, response)
			if err == nil || r >= retries {
				return err
			}
		default:
			return e.NewCustomHTTPStatus(e.StatusBadRequest, "", "invalid requestType")
		}

		baseDelay = time.Duration(int64(baseDelay) * int64(r+1))

		log.Printf("Attempt %d failed; retrying in %v", r+1, baseDelay)
		select {
		case <-time.After(baseDelay):
		case <-ctx.Done():
			return err
		}
	}
}

// HandleMultipartWriter handle the multipart requests and write parts to the defined path
func (rs *restService) HandleMultipartWriter(request *Api, pathToWrite string, response interface{}) e.IError {
	if err := os.MkdirAll(pathToWrite, os.ModePerm); err != nil {
		return e.NewCustomHTTPStatus(e.StatusInternalServerError, "", err.Error())
	}

	urlPath := path.Join(rs.url.Path, fmt.Sprintf(request.endpoint, request.arguments...))
	endpoint, err := rs.url.Parse(urlPath)
	if err != nil {
		return e.NewCustomHTTPStatus(e.StatusInternalServerError, "", err.Error())
	}

	req := &http.Request{
		Method: request.method,
		URL:    endpoint,
		Header: rs.headers,
	}

	resp, err := rs.Do(req)
	if err != nil {
		return e.NewCustomHTTPStatus(e.StatusServiceUnavailable, "", err.Error())
	}
	defer resp.Body.Close()

	// Check if the response is successful
	if resp.StatusCode != http.StatusOK {
		return e.NewCustomHTTPStatus(e.StatusCode(resp.StatusCode))
	}

	// Check if the response is multipart
	contentType := resp.Header.Get("Content-Type")
	if !isMultipart(contentType) {
		return e.NewCustomHTTPStatus(e.StatusBadRequest, "", "content-type is not multipart")
	}

	// Create a multipart reader
	multipartReader := multipart.NewReader(resp.Body, boundaryFromContentType(contentType))
	// Read each part
	for {
		part, err := multipartReader.NextPart()
		if err != nil {
			break
		}
		defer part.Close()

		// Create the file on the client server
		clientFilePath := filepath.Join(pathToWrite, part.FileName())
		file, err := os.Create(clientFilePath)
		if err != nil {
			return e.NewCustomHTTPStatus(e.StatusInternalServerError, "", err.Error())
		}
		defer file.Close()

		// Copy the part content to the client file
		_, err = io.Copy(file, part)
		if err != nil {
			return e.NewCustomHTTPStatus(e.StatusInternalServerError, "", err.Error())
		}

		log.Printf("file %s created successfully on the client server\n", part.FileName())
	}
	return nil
}

// Check if the content type is multipart
func isMultipart(contentType string) bool {
	return strings.HasPrefix(contentType, "multipart/form-data")
}

// Extract the boundary from the content type
func boundaryFromContentType(contentType string) string {
	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return ""
	}
	return params["boundary"]
}

// handleRequest sends a HTTP(S), placing results into the response object
func (rs *restService) HandleRequest(request *Api, response interface{}) e.IError {
	urlPath := path.Join(rs.url.Path, fmt.Sprintf(request.endpoint, request.arguments...))
	endpoint, err := rs.url.Parse(urlPath)
	if err != nil {
		return e.NewCustomHTTPStatus(e.StatusInternalServerError, "", err.Error())
	}

	var readCloser io.ReadCloser
	if request.body != nil {
		outbuf, err := json.Marshal(request.body)
		if err != nil {
			return e.NewCustomHTTPStatus(e.StatusInternalServerError, "", err.Error())
		}
		readCloser = ioutil.NopCloser(bytes.NewBuffer(outbuf))
	}

	req := &http.Request{
		Method: request.method,
		URL:    endpoint,
		Body:   readCloser,
		Header: rs.headers,
	}

	resp, err := rs.Do(req)
	if err != nil {
		return e.NewCustomHTTPStatus(e.StatusServiceUnavailable, "", err.Error())
	}

	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		if err = json.NewDecoder(resp.Body).Decode(response); err != nil {
			return e.NewCustomHTTPStatus(e.StatusInternalServerError, "", err.Error())
		}
		return nil
	}

	var failure RestError
	failure.Code = 200
	if err := json.NewDecoder(resp.Body).Decode(&failure); err != nil {
		return e.NewCustomHTTPStatus(e.StatusInternalServerError, "", err.Error())
	}

	return nil
}
