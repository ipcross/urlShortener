package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ipcross/urlShortener/internal/config"
	"github.com/ipcross/urlShortener/internal/handlers"
	"github.com/ipcross/urlShortener/internal/repository"
	"github.com/ipcross/urlShortener/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type want struct {
	code        int
	url         string
	contentType string
}

func TestHandlers(t *testing.T) {
	cfg := config.GetConfig()
	store := repository.NewStore()
	mapperService := service.NewMapper(cfg, store)
	h := handlers.NewHandlers(mapperService, cfg)
	tests := []struct {
		name string
		want want
	}{
		{
			name: "web post",
			want: want{
				code:        201,
				url:         `/`,
				contentType: "text/plain",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			longURL := strings.NewReader("https://yandex.ru")
			request := httptest.NewRequest(http.MethodPost, test.want.url, longURL)
			w := httptest.NewRecorder()
			h.PostHandler(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer dclose(res.Body)
			_, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}

	t.Run("API shorter", func(t *testing.T) {
		postBody := map[string]interface{}{
			"url": "www.ru",
		}
		body, _ := json.Marshal(postBody)
		request := httptest.NewRequest(http.MethodGet, "/api/shorter", bytes.NewReader(body))
		w := httptest.NewRecorder()
		h.APIHandler(w, request)

		res := w.Result()
		var m map[string]interface{}
		err := json.NewDecoder(res.Body).Decode(&m)
		require.NoError(t, err)
		assert.Contains(t, m, "result")
		assert.Equal(t, 201, res.StatusCode)
		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
		defer dclose(res.Body)
	})

	t.Run("Not found", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/bad_hash", http.NoBody)
		w := httptest.NewRecorder()
		h.GetHandler(w, request)

		res := w.Result()
		assert.Equal(t, 400, res.StatusCode)
		defer dclose(res.Body)
	})

	t.Run("Bad request", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPut, "/", http.NoBody)
		w := httptest.NewRecorder()
		h.BadRequestHandler(w, request)

		res := w.Result()
		assert.Equal(t, 400, res.StatusCode)
		defer dclose(res.Body)
	})
}

func dclose(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Fatal(err)
	}
}
