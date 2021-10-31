package handlers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/xsoroton/autopilot-proxy/store"
)

func getTestHandlers(remoteHost string) (handlers Handlers) {
	var err error
	handlers.Remote, err = url.Parse(remoteHost)
	if err != nil {
		panic(err)
	}
	// run memory implementation of store for testing
	handlers.CacheStore = store.NewMemStore()
	handlers.Proxy = httputil.NewSingleHostReverseProxy(handlers.Remote)
	return
}

func TestGetContact(t *testing.T) {
	Convey("Test GET Contact", t, func() {
		// Mock Remote API
		respValue := `{}`
		mockAPI := MockHTTPServer(respValue, http.StatusOK)
		h := getTestHandlers(mockAPI.URL)

		// Run Proxy
		gin.SetMode(gin.TestMode)
		r := gin.New()
		r.GET("/v1/contact/*contact_id", h.GetContact)
		host := ":9980"
		go func() {
			r.Run(host)
		}()
		// wait for proxy just in case
		time.Sleep(time.Millisecond * 100)

		var client http.Client
		req, err := http.NewRequest(http.MethodGet, "http://0.0.0.0:9980/v1/contact/test", nil)
		So(err, ShouldBeNil)

		res, err := client.Do(req)
		So(err, ShouldBeNil)
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		So(err, ShouldBeNil)
		So(res.StatusCode, ShouldEqual, http.StatusOK)
		So(string(body), ShouldEqual, respValue)

		Convey("Test cache is set", func() {
			v, err := h.CacheStore.Get("test")
			So(err, ShouldBeNil)
			So(string(v), ShouldEqual, respValue)
		})
	})
}

func TestPostContact(t *testing.T) {
	Convey("Test POST Contact", t, func() {
		// Mock Remote API
		clientID := "person_BEDEF3B9-8B84-4F5F-AA58-22D025DDA683"
		respValue := fmt.Sprintf(`{"contact_id": "%s"}`, clientID)
		mockAPI := MockHTTPServer(respValue, http.StatusOK)
		h := getTestHandlers(mockAPI.URL)

		// Run Proxy
		gin.SetMode(gin.TestMode)
		r := gin.New()
		r.POST("/v1/contact", h.PostContact)
		host := ":9990"
		go func() {
			r.Run(host)
		}()
		// wait for proxy just in case
		time.Sleep(time.Millisecond * 100)

		// set cache, late check that POST invalidated
		err := h.CacheStore.Set(clientID, []byte(`{}`))
		So(err, ShouldBeNil)

		var client http.Client
		req, err := http.NewRequest(http.MethodPost, "http://0.0.0.0:9990/v1/contact", bytes.NewBuffer([]byte(`{}`)))
		So(err, ShouldBeNil)

		res, err := client.Do(req)
		So(err, ShouldBeNil)
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		So(err, ShouldBeNil)
		So(res.StatusCode, ShouldEqual, http.StatusOK)
		So(string(body), ShouldEqual, respValue)

		Convey("Test POST invalidate cache", func() {
			_, err := h.CacheStore.Get(clientID)
			So(err, ShouldBeError)
		})
	})
}

// MockHTTPServer mock http server response and responseCode
func MockHTTPServer(response string, responseCode int) *httptest.Server {
	handler := handler(response, responseCode)
	return httptest.NewServer(http.HandlerFunc(handler))
}

func handler(response string, responseCode int) (handler func(w http.ResponseWriter, req *http.Request)) {
	handler = func(w http.ResponseWriter, req *http.Request) {
		// mock response code
		w.WriteHeader(responseCode)
		// mock response body
		w.Write([]byte(response))
	}
	return
}
