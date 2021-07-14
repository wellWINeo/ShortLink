package handler

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/skip2/go-qrcode"
	"github.com/wellWINeo/ShortLink/pkg/service"
	"golang.org/x/net/html"
)

func TestHandler_createLink(t *testing.T) {
	mock := service.NewLinksMock(t)
	h := NewHandler(&service.Service{Links: mock}, "website", "localhost:8000")

	testTable := []struct {
		name          string
		fieldName     string
		input         string
		expect        string
		expectStatus  int
		mockBehaviour func(string)
	}{
		// TODO fill 'expect' fields
		{
			name:         "OK",
			fieldName:    "origin_link",
			input:        "https://example.com",
			expect:       "http://exmaple.com",
			expectStatus: http.StatusMovedPermanently,
			mockBehaviour: func(input string) {
				mock.CreateLinkMock.
					Expect(input).
					Return("someTestHash", nil)
			},
		},
		{
			name:         "Wrong field name",
			fieldName:    "wrong_name",
			input:        "https://example.com",
			expect:       "",
			expectStatus: http.StatusBadRequest,
			mockBehaviour: func(input string) {
				mock.CreateLinkMock.
					Expect("").
					Return("someTestHash", nil)
			},
		},
		{
			name:         "Empty field value",
			fieldName:    "origin_link",
			input:        "",
			expect:       "",
			expectStatus: http.StatusBadRequest,
			mockBehaviour: func(input string) {
				mock.CreateLinkMock.
					Expect(input).
					Return("someTestHash", nil)
			},
		},
		{
			name:         "Validation failed",
			fieldName:    "origin_link",
			input:        "https://example.com",
			expect:       "",
			expectStatus: http.StatusBadRequest,
			mockBehaviour: func(input string) {
				mock.CreateLinkMock.
					Expect(input).
					Return("", errors.New("Wrong url"))
			},
		},
	}

	for _, testCase := range testTable {

		testCase.mockBehaviour(testCase.input)

		formData := url.Values{}
		formData.Add(testCase.fieldName, testCase.input)

		router := gin.Default()
		router.POST("/create", h.createLink)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/create",
			strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		router.ServeHTTP(w, req)

		// fmt.Printf("[%s]\nExpect: %d\nGot: %d\n",
		// 	testCase.name, testCase.expectStatus, w.Code)

		assert.Equal(t, testCase.expectStatus, w.Code)

		// TODO add body assertion
		// assert.Equal(t, testCase.expect, rec.Body.String())
	}
}

func TestHandler_getLink(t *testing.T) {
	mock := service.NewLinksMock(t)
	h := NewHandler(&service.Service{Links: mock}, "website", "localhost:8000")

	testTable := []struct {
		name          string
		input         string
		expect        string
		expectStatus  int
		mockBehaviour func(string)
	}{
		{
			name:         "OK",
			input:        "id123",
			expect:       "http://example.com",
			expectStatus: http.StatusMovedPermanently,
			mockBehaviour: func(input string) {
				mock.GetLinkMock.
					Expect(input).
					Return("http://example.com", nil)
			},
		},
		{
			name:         "Empty",
			input:        "",
			expect:       "",
			expectStatus: http.StatusNotFound,
			mockBehaviour: func(input string) {
				mock.GetLinkMock.
					Expect(input).
					Return("", nil)
			},
		},
		{
			name:         "Not exist",
			input:        "no123exist",
			expect:       "",
			expectStatus: http.StatusNotFound,
			mockBehaviour: func(input string) {
				mock.GetLinkMock.
					Expect(input).
					Return("", errors.New("No such record"))
			},
		},
	}

	for _, testCase := range testTable {

		testCase.mockBehaviour(testCase.input)

		router := gin.Default()
		router.GET("/:url", h.getLink)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/" + testCase.input, nil)

		router.ServeHTTP(w, req)

		// fmt.Printf("[%s]\nExpect: %d\nGot: %d\n",
		// 	testCase.name, testCase.expectStatus, w.Code)

		assert.Equal(t, testCase.expectStatus, w.Code)

		// TODO add body assertion
		// assert.Equal(t, testCase.expect, rec.Body.String())
	}
}

func extractTagAttr(body, tag, attr string) ([]string, error) {
	var values []string
	reader := strings.NewReader(body)
	tokenizer := html.NewTokenizer(reader)

	for {
		tok := tokenizer.Next()
		if tok == html.ErrorToken {
			if tokenizer.Err() == io.EOF {
				return values, nil
			}
			fmt.Printf("Error token parsing: %s", tokenizer.Err())
			return values, errors.New(tokenizer.Err().Error())
		}

		tagName, ok := tokenizer.TagName()
		if  ok && string(tagName) == tag {
			for {
				key, value, more := tokenizer.TagAttr()
				if string(key) == attr {
					values = append(values, string(value))
				}
				if !more {
					break
				}
			}
		}
	}
}

func TestHandler_getQR(t *testing.T) {
	mock := service.NewLinksMock(t)
	h := NewHandler(&service.Service{Links: mock}, "../../website", "localhost:8000")

	testTable := []struct {
		name          string
		input         string
		expect        string
		expectStatus  int
		mockBehaviour func(string)
	}{
		{
			name:         "OK",
			input:        "id123",
			expect:       "http://example.com",
			expectStatus: http.StatusOK,
			mockBehaviour: func(input string) {
				img, _ := qrcode.Encode("http://example.com", qrcode.Medium, 256)
				mock.GetQRMock.
					Expect(input).
					Return(img, "http://example.com", nil)
			},
		},
	}

	for _, testCase := range testTable {

		testCase.mockBehaviour(testCase.input)

		router := gin.Default()
		router.GET("/qr/:url", h.getQR)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/qr/" + testCase.input, nil)

		router.ServeHTTP(w, req)

		var value string
		imgs, err := extractTagAttr(w.Body.String(), "img", "src")
		if len(imgs) != 1 {
			t.Fail()
			continue
		}
		val := strings.Split(string(imgs[0]), " ")
		if len(val) != 0 {
			value = val[len(val)-1]
		}
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			t.Fail()
		}
		fmt.Printf("Img: %s\n", value)
		qr, _ := qrcode.Encode("http://example.com", qrcode.Medium, 256)
		str := base64.StdEncoding.EncodeToString(qr)

		// assertition
		assert.Equal(t, testCase.expectStatus, w.Code)
		assert.Equal(t, str, value)
	}
}
