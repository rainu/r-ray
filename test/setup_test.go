package test

import (
	"context"
	"fmt"
	ihttp "github.com/rainu/r-ray/internal/http"
	"github.com/rainu/r-ray/internal/http/controller"
	"github.com/rainu/r-ray/internal/processor"
	"github.com/rainu/r-ray/internal/store"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
)

const (
	username = "user"
	password = "pw"

	headerPrefix  = "R-"
	fwdReqHeader  = headerPrefix + "Forward-Request"
	fwdRespHeader = headerPrefix + "Forward-Response"
)

var (
	testServer *httptest.Server
	mock       func(w http.ResponseWriter, r *http.Request)
	emptyMock  = func(w http.ResponseWriter, r *http.Request) {}
	appPort    int
	shutdowner interface{ Shutdown(context.Context) error }
)

func init() {
	testServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if mock != nil {
			//call mock
			mock(w, r)
		}
	}))
	appPort = freePort()

	userStore := store.NewUser()
	userStore.Add(username, password)

	p := processor.New(userStore)
	h := controller.NewProxy(headerPrefix, fwdReqHeader, fwdRespHeader, p)
	toTest := ihttp.NewServer(fmt.Sprintf(":%d", appPort), h)
	shutdowner = toTest

	//start test server
	go toTest.ListenAndServe()
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

func rrayUrl(path string) string {
	return fmt.Sprintf(
		"http://localhost:%d/?url=%s",
		appPort,
		url.QueryEscape(testServer.URL+path),
	)
}

func validRequest(method, url string, header http.Header, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, rrayUrl(url), body)
	if err != nil {
		panic(err)
	}

	req.SetBasicAuth(username, password)
	if header != nil {
		req.Header = header
	}

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
