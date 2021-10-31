package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"

	"github.com/xsoroton/autopilot-proxy/models"
)

func (h *Handlers) GetContact(c *gin.Context) {
	contextID := c.Param("contact_id")[1:]
	val, err := h.CacheStore.Get(contextID)
	if err == nil && len(val) != 0 {
		c.String(http.StatusOK, string(val))
		return
	}

	h.Proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Header.Del("Accept-Encoding") // todo: add handle of zip context in future
		req.Host = h.Remote.Host
		req.URL.Scheme = h.Remote.Scheme
		req.URL.Host = h.Remote.Host
		req.URL.Path = fmt.Sprintf("/v1/contact/%s", contextID)
	}
	//proxy to read response body and cache it
	h.Proxy.ModifyResponse = func(resp *http.Response) (err error) {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		err = h.CacheStore.Set(contextID, b)
		if err != nil {
			return err
		}
		//logrus.Infof("%s", b)
		resp.Body = ioutil.NopCloser(bytes.NewReader(b))
		return
	}
	h.Proxy.ServeHTTP(c.Writer, c.Request)
}

func (h *Handlers) PostContact(c *gin.Context) {
	h.Proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Header.Del("Accept-Encoding") // todo: add handle of zip context in future
		req.Host = h.Remote.Host
		req.URL.Scheme = h.Remote.Scheme
		req.URL.Host = h.Remote.Host
		req.URL.Path = "/v1/contact"
	}
	// proxy to read response body and invalidate cache
	h.Proxy.ModifyResponse = func(resp *http.Response) (err error) {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logrus.Error(err)
			return err
		}
		defer resp.Body.Close()
		var postContactResponse models.PostContactResponse
		err = json.Unmarshal(b, &postContactResponse)
		if err != nil {
			return err
		}
		err = h.CacheStore.Remove(postContactResponse.ContactID)
		if err != nil {
			return err
		}
		//logrus.Infof("%s", b)
		resp.Body = ioutil.NopCloser(bytes.NewReader(b))
		return
	}
	h.Proxy.ServeHTTP(c.Writer, c.Request)
}
