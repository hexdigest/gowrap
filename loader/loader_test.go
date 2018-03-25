package loader

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gojuno/minimock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	l := New(nil)
	assert.Equal(t, Loader{client: http.DefaultClient}, l)
}

func TestLoader_Load(t *testing.T) {
	clientError := errors.New("client error")

	tests := []struct {
		name    string
		init    func(t minimock.Tester) Loader
		inspect func(r Loader, t *testing.T) //inspects Loader after execution of Load

		path string

		want1      []byte
		want2      string
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation
	}{
		{
			name: "HTTPS URL",
			path: "https://host/template",
			init: func(t minimock.Tester) Loader {
				return Loader{client: newHTTPClientMock(t).DoMock.Return(nil, clientError)}
			},
			want2:   "https://host/template",
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, clientError, err)
			},
		},
		{
			name: "github reference",
			path: "github/template",
			init: func(t minimock.Tester) Loader {
				return Loader{client: newHTTPClientMock(t).DoMock.Return(nil, clientError)}
			},
			want2:   "",
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, clientError, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)

			receiver := tt.init(mc)

			got1, got2, err := receiver.Load(tt.path)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			assert.Equal(t, tt.want1, got1, "Loader.Load returned unexpected result")

			assert.Equal(t, tt.want2, got2, "Loader.Load returned unexpected result")

			if tt.wantErr {
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestLoader_List(t *testing.T) {
	clientError := errors.New("client error")
	tests := []struct {
		name    string
		init    func(t minimock.Tester) Loader
		inspect func(r Loader, t *testing.T) //inspects Loader after execution of List

		want1      []string
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation
	}{
		{
			name: "client error",
			init: func(t minimock.Tester) Loader {
				return Loader{client: newHTTPClientMock(t).DoMock.Return(nil, clientError)}
			},
			want1:   nil,
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, clientError, err)
			},
		},
		{
			name: "unmarshal error",
			init: func(t minimock.Tester) Loader {
				r := &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(strings.NewReader("{"))}
				return Loader{client: newHTTPClientMock(t).DoMock.Return(r, nil)}
			},
			want1:   nil,
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				_, ok := err.(*json.SyntaxError)
				assert.True(t, ok)
			},
		},
		{
			name: "success",
			init: func(t minimock.Tester) Loader {
				r := &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(strings.NewReader(`{
					"tree": [
						{
							"path": "templates/opentracing"
						}
					]
				}`))}
				return Loader{client: newHTTPClientMock(t).DoMock.Return(r, nil)}
			},
			want1: []string{"opentracing"},
			inspectErr: func(err error, t *testing.T) {
				_, ok := err.(*json.SyntaxError)
				assert.True(t, ok)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)

			receiver := tt.init(mc)

			got1, err := receiver.List()

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			assert.Equal(t, tt.want1, got1, "Loader.List returned unexpected result")

			if tt.wantErr {
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestLoader_get(t *testing.T) {
	errBodyClose := errors.New("body close error")

	tests := []struct {
		name    string
		init    func(t minimock.Tester) Loader
		inspect func(r Loader, t *testing.T) //inspects Loader after execution of get

		url string

		want1      []byte
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation
	}{
		{
			name:    "invalid URL",
			init:    func(t minimock.Tester) Loader { return Loader{} },
			url:     "https://\\",
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				_, ok := err.(*url.Error)
				assert.True(t, ok)
			},
		},

		{
			name: "unexpected status code",
			init: func(t minimock.Tester) Loader {
				r := &http.Response{StatusCode: http.StatusInternalServerError, Body: ioutil.NopCloser(strings.NewReader(""))}
				return Loader{client: newHTTPClientMock(t).DoMock.Return(r, nil)}
			},
			url:     "https://host",
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, errUnexpectedStatusCode, errors.Cause(err))
			},
		},

		{
			name: "body close error",
			init: func(t minimock.Tester) Loader {
				rcMock := NewReadCloserMock(t)
				rcMock.ReadFunc = func(p []byte) (int, error) {
					return 0, io.EOF
				}
				rcMock.CloseMock.Return(errBodyClose)

				r := &http.Response{StatusCode: http.StatusOK, Body: rcMock}
				return Loader{client: newHTTPClientMock(t).DoMock.Return(r, nil)}
			},
			url:     "https://host",
			wantErr: true,
			want1:   []byte{},
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, errBodyClose, errors.Cause(err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Finish()

			receiver := tt.init(mc)

			got1, err := receiver.get(tt.url)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			assert.Equal(t, tt.want1, got1, "Loader.get returned unexpected result")

			if tt.wantErr {
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestLoader_fetchFromGithub(t *testing.T) {
	clientError := errors.New("client error")

	tests := []struct {
		name    string
		init    func(t minimock.Tester) Loader
		inspect func(r Loader, t *testing.T) //inspects Loader after execution of fetchFromGithub

		templateName string

		want1      []byte
		want2      string
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation
	}{
		{
			name: "client error",
			init: func(t minimock.Tester) Loader {
				return Loader{client: newHTTPClientMock(t).DoMock.Return(nil, clientError)}
			},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, clientError, err)
			},
		},
		{
			name: "unmarshal error",
			init: func(t minimock.Tester) Loader {
				r := &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(strings.NewReader("{"))}
				return Loader{client: newHTTPClientMock(t).DoMock.Return(r, nil)}
			},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				_, ok := errors.Cause(err).(*json.SyntaxError)
				assert.True(t, ok)
			},
		},
		{
			name: "no commits",
			init: func(t minimock.Tester) Loader {
				r := &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(strings.NewReader("[]"))}
				return Loader{client: newHTTPClientMock(t).DoMock.Return(r, nil)}
			},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, errTemplateNotFound, err)
			},
		},
		{
			name: "template fetch error",
			init: func(t minimock.Tester) Loader {
				var i int
				client := newHTTPClientMock(t)
				client.DoFunc = func(r *http.Request) (*http.Response, error) {
					if i == 0 {
						i++
						return &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(strings.NewReader(`[{"SHA":"hash"}]`))}, nil
					}
					return nil, clientError
				}

				return Loader{client: client}
			},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, clientError, errors.Cause(err))
			},
		},
		{
			name:         "success",
			templateName: "opentracing",
			init: func(t minimock.Tester) Loader {
				var i int
				client := newHTTPClientMock(t)
				client.DoFunc = func(r *http.Request) (*http.Response, error) {
					if i == 0 {
						i++
						return &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(strings.NewReader(`[{"SHA":"hash"}]`))}, nil
					}

					return &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(strings.NewReader(`template body`))}, nil
				}

				return Loader{client: client}
			},
			want1: []byte("template body"),
			want2: "https://raw.githubusercontent.com/hexdigest/gounit/hash/templates/opentracing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)

			receiver := tt.init(mc)

			got1, got2, err := receiver.fetchFromGithub(tt.templateName)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			assert.Equal(t, tt.want1, got1, "Loader.fetchFromGithub returned unexpected result")

			assert.Equal(t, tt.want2, got2, "Loader.fetchFromGithub returned unexpected result")

			if tt.wantErr {
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err)
			}

		})
	}
}
