package test

import (
	"context"
	"fmt"
	ihttp "github.com/rainu/r-ray/internal/http"
	"github.com/rainu/r-ray/internal/http/controller"
	"github.com/rainu/r-ray/internal/processor"
	"github.com/rainu/r-ray/internal/store"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"
)

const (
	username = "user"
	password = "pw"

	headerPrefix        = "R-"
	fwdReqHeaderPrefix  = headerPrefix + "Forward-Request-"
	fwdRespHeaderPrefix = headerPrefix + "Forward-Response-"
)

var (
	testServer *httptest.Server
	mock       func(w http.ResponseWriter, r *http.Request)
	appPort    int
	appBaseUrl string
	shutdowner interface{ Shutdown(context.Context) error }
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)

	testServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if mock != nil {
			//call mock
			mock(w, r)
		}
	}))
	appPort = freePort()
	appBaseUrl = fmt.Sprintf("http://localhost:%d/", appPort)

	userStore := store.NewUser()
	userStore.Add(username, password)

	p := processor.New(userStore)
	h := controller.NewProxy(headerPrefix, fwdReqHeaderPrefix, fwdRespHeaderPrefix, p)
	toTest := ihttp.NewServer(fmt.Sprintf(":%d", appPort), h)
	shutdowner = toTest

	//start test server
	go toTest.ListenAndServe()

	//give some time to startup the proxy
	time.Sleep(50 * time.Millisecond)
}

func freePort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	return listener.Addr().(*net.TCPAddr).Port
}

func appUrl(path string) string {
	return fmt.Sprintf(
		"%s%s",
		appBaseUrl,
		url.QueryEscape(testServer.URL+path),
	)
}

func validRequest(method, url string, header http.Header, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, appUrl(url), body)
	if err != nil {
		panic(err)
	}

	if header != nil {
		req.Header = header
	}
	req.SetBasicAuth(username, password)

	return req
}

func do(req *http.Request) *http.Response {
	resp, err := testServer.Client().Do(req)
	if err != nil {
		panic(err)
	}

	//remove the following header because they are badly testable
	resp.Header.Del("Date")
	resp.Header.Del(headerPrefix + "Date")

	return resp
}
