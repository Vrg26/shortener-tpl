package shorturl

import (
	"bytes"
	"encoding/json"
	"github.com/Vrg26/shortener-tpl/internal/app/shorturl/db"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}

func Test_handler_AddUrl(t *testing.T) {
	st := db.NewMemoryStorage()
	s := NewService(st)
	handlerSU := NewHandler(*s, "http://localhost")
	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name    string
		request string
		body    string
		want    want
	}{
		{
			name:    "success test",
			request: "/",
			body:    "https://twitter.com",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  201,
			},
		},
		{
			name:    "should return error bad request. Empty body",
			request: "/",
			body:    "",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
			},
		},
		{
			name:    "should return not found. Invalid path",
			request: "/test",
			body:    "testestset",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  404,
			},
		},
		{
			name:    "should return error bad request. Invalid URL",
			request: "/",
			body:    "testestset",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handlerSU.AddTextURL)
			h.ServeHTTP(w, request)
			res := w.Result()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))

			if tt.want.statusCode == http.StatusCreated {
				result, err := ioutil.ReadAll(res.Body)
				require.NoError(t, err)
				err = res.Body.Close()
				require.NoError(t, err)
				_, err = url.ParseRequestURI(string(result))
				assert.NoError(t, err)
			}
		})
	}
}

func Test_handler_AddShorten(t *testing.T) {
	st := db.NewMemoryStorage()
	s := NewService(st)
	handlerSU := NewHandler(*s, "http://localhost")
	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name    string
		request string
		body    string
		want    want
	}{
		{
			name:    "success test",
			request: "/api/shorten",
			body:    `{ "url":"https://twitter.com"}`,
			want: want{
				contentType: "application/json; charset=utf-8",
				statusCode:  201,
			},
		},
		{
			name:    "should return error bad request. Empty body",
			request: "/api/shorten",
			body:    "",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
			},
		},
		{
			name:    "should return error bad request. Invalid URL",
			request: "/api/shorten",
			body:    "testestset",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handlerSU.AddJsonURL)
			h.ServeHTTP(w, request)
			res := w.Result()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))

			if tt.want.statusCode == http.StatusCreated {

				var result ResponseURL
				err := json.NewDecoder(res.Body).Decode(&result)
				require.NoError(t, err)
				err = res.Body.Close()
				require.NoError(t, err)
				_, err = url.ParseRequestURI(result.Result)
				assert.NoError(t, err)
			}
		})
	}
}

func Test_handler_GetUrl(t *testing.T) {
	r := chi.NewRouter()
	st := db.NewMemoryStorage()
	s := NewService(st)
	h := NewHandler(*s, "")
	h.Register(r)

	idURL, err := st.Add("https://practicum.yandex.ru")
	require.NoError(t, err)

	ts := httptest.NewServer(r)
	defer ts.Close()

	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "success test",
			request: "/" + idURL,
			want: want{
				contentType: "text/html; charset=utf-8",
				statusCode:  200,
			},
		},
		{
			name:    "should return not found",
			request: "/1234",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  404,
			},
		},
		{
			name:    "should return method not allowed",
			request: "/",
			want: want{
				contentType: "",
				statusCode:  405,
			},
		},
	}
	defer ts.Close()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, _ := testRequest(t, ts, http.MethodGet, tt.request)
			defer res.Body.Close()
			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
