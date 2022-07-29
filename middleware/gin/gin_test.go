package gin_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	mmetrics "github.com/yangshun2005/go-prometheus-http/internal/mocks/metrics"
	"github.com/yangshun2005/go-prometheus-http/metrics"
	"github.com/yangshun2005/go-prometheus-http/middleware"
	ginmiddleware "github.com/yangshun2005/go-prometheus-http/middleware/gin"
)

func TestMiddleware(t *testing.T) {
	tests := map[string]struct {
		handlerID   string
		config      middleware.Config
		req         func() *http.Request
		mock        func(m *mmetrics.Recorder)
		handler     func() gin.HandlerFunc
		expRespCode int
		expRespBody string
	}{
		"A default HTTP middleware should call the recorder to measure.": {
			req: func() *http.Request {
				return httptest.NewRequest(http.MethodPost, "/test", nil)
			},
			mock: func(m *mmetrics.Recorder) {
				expHTTPReqProps := metrics.HTTPReqProperties{
					ID:      "/test",
					Service: "",
					Method:  "POST",
					Code:    "202",
				}
				m.On("ObserveHTTPRequestDuration", mock.Anything, expHTTPReqProps, mock.Anything).Once()
				m.On("ObserveHTTPResponseSize", mock.Anything, expHTTPReqProps, int64(5)).Once()

				expHTTPProps := metrics.HTTPProperties{
					ID:      "/test",
					Service: "",
				}
				m.On("AddInflightRequests", mock.Anything, expHTTPProps, 1).Once()
				m.On("AddInflightRequests", mock.Anything, expHTTPProps, -1).Once()
			},
			handler: func() gin.HandlerFunc {
				return func(c *gin.Context) {
					c.String(202, "test1")
				}
			},
			expRespCode: 202,
			expRespBody: "test1",
		},

		"A default HTTP middleware using JSON should call the recorder to measure (Regression test: https://github.com/yangshun2005/go-prometheus-http/issues/31).": {
			req: func() *http.Request {
				return httptest.NewRequest(http.MethodPost, "/test", nil)
			},
			mock: func(m *mmetrics.Recorder) {
				expHTTPReqProps := metrics.HTTPReqProperties{
					ID:      "/test",
					Service: "",
					Method:  "POST",
					Code:    "202",
				}
				m.On("ObserveHTTPRequestDuration", mock.Anything, expHTTPReqProps, mock.Anything).Once()
				m.On("ObserveHTTPResponseSize", mock.Anything, expHTTPReqProps, int64(14)).Once()

				expHTTPProps := metrics.HTTPProperties{
					ID:      "/test",
					Service: "",
				}
				m.On("AddInflightRequests", mock.Anything, expHTTPProps, 1).Once()
				m.On("AddInflightRequests", mock.Anything, expHTTPProps, -1).Once()
			},
			handler: func() gin.HandlerFunc {
				return func(c *gin.Context) {
					c.JSON(202, map[string]string{"test": "one"})
				}
			},
			expRespCode: 202,
			expRespBody: `{"test":"one"}`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			// Mocks.
			mr := &mmetrics.Recorder{}
			test.mock(mr)

			// Create our instance with the middleware.
			mdlw := middleware.New(middleware.Config{Recorder: mr})
			engine := gin.New()
			req := test.req()
			engine.Handle(req.Method, req.URL.Path,
				ginmiddleware.Handler(test.handlerID, mdlw),
				test.handler())

			// Make the request.
			resp := httptest.NewRecorder()
			engine.ServeHTTP(resp, req)

			// Check.
			mr.AssertExpectations(t)
			assert.Equal(test.expRespCode, resp.Result().StatusCode)
			gotBody, err := ioutil.ReadAll(resp.Result().Body)
			require.NoError(err)
			assert.Equal(test.expRespBody, string(gotBody))
		})
	}
}
