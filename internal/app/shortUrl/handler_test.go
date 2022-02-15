package shortUrl

import (
	"bytes"
	"github.com/Vrg26/shortener-tpl/internal/app/shortUrl/db"
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
	handlerSU := NewHandler(*s)
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
			name:    "Correct request",
			request: "/",
			body:    "https://twitter.com",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  201,
			},
		},
		{
			name:    "Empty body",
			request: "/",
			body:    "",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
			},
		},
		{
			name:    "Incorrect path",
			request: "/test",
			body:    "testestset",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  404,
			},
		},
		{
			name:    "Incorrect url",
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
			h := http.HandlerFunc(handlerSU.AddUrl)
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

func Test_handler_GetUrl(t *testing.T) {
	r := chi.NewRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()
	st := db.NewMemoryStorage()
	idUrl, err := st.Add("https://jsonplaceholder.typicode.com/posts")
	require.NoError(t, err)
	s := NewService(st)
	handlerSU := NewHandler(*s)
	handlerSU.RegisterChi(r)
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
			name:    "correct test",
			request: "/" + idUrl,
			want: want{
				contentType: "text/html; charset=utf-8",
				statusCode:  307,
			},
		},
		{
			name:    "url not found",
			request: "/1234",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  404,
			},
		},
		{
			name:    "empty get request",
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
			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}