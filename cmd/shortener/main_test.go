package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMyHandler(t *testing.T) {
	type want struct {
		code           int
		response       string
		contentType    string
		headerLocation string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "post",
			want: want{
				code:        201,
				response:    `/1`,
				contentType: "text/plain",
			},
		},
		{
			name: "get",
			want: want{
				code:           307,
				headerLocation: "http://yandex.ru",
			},
		},
		{
			name: "bad_request",
			want: want{
				code: 400,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			switch test.name {
			case "post":
				longURL := strings.NewReader("http://yandex.ru")
				request := httptest.NewRequest(http.MethodPost, "/", longURL)
				w := httptest.NewRecorder()
				PostHandler(w, request)

				res := w.Result()
				assert.Equal(t, test.want.code, res.StatusCode)

				defer res.Body.Close()
				resBody, err := io.ReadAll(res.Body)

				require.NoError(t, err)
				assert.Equal(t, test.want.response, string(resBody))
				assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			case "get":
				request := httptest.NewRequest(http.MethodGet, "/1", nil)
				w := httptest.NewRecorder()
				GetHandler(w, request)

				res := w.Result()
				assert.Equal(t, test.want.code, res.StatusCode)
				assert.Equal(t, test.want.headerLocation, res.Header.Get("Location"))
				defer res.Body.Close()
			default:
				request := httptest.NewRequest(http.MethodPut, "/", nil)
				w := httptest.NewRecorder()
				BadRequestHandler(w, request)

				res := w.Result()
				assert.Equal(t, test.want.code, res.StatusCode)
				defer res.Body.Close()
			}
		})
	}
}
