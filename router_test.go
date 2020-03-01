package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter_Method(t *testing.T) {

	router := NewRouter()
	testMethods := []struct {
		route  *Route
		method string
		expect int
	}{
		{route: router.GET("/a", testHandler), method: http.MethodGet, expect: http.StatusOK},
		{route: router.POST("/b", testHandler), method: http.MethodPost, expect: http.StatusOK},
		{route: router.DELETE("/c", testHandler), method: http.MethodDelete, expect: http.StatusOK},
		{route: router.PUT("/e", testHandler), method: http.MethodPut, expect: http.StatusOK},
		{route: router.PATCH("/f", testHandler), method: http.MethodPatch, expect: http.StatusOK},

		{route: router.GET("/a", testHandler), method: http.MethodPost, expect: http.StatusMethodNotAllowed},
		{route: router.POST("/b", testHandler), method: http.MethodGet, expect: http.StatusMethodNotAllowed},
		{route: router.DELETE("/c", testHandler), method: http.MethodPost, expect: http.StatusMethodNotAllowed},
		{route: router.PUT("/e", testHandler), method: http.MethodPatch, expect: http.StatusMethodNotAllowed},
		{route: router.PATCH("/f", testHandler), method: http.MethodDelete, expect: http.StatusMethodNotAllowed},
	}

	for _, tm := range testMethods {
		router.AddRoute(tm.route.path, tm.method, testHandler)
		req := httptest.NewRequest(tm.method, tm.route.path, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != tm.expect {
			t.Errorf("Expected status code %d (got %d) on method %s", tm.expect, w.Code, tm.method)
		}
	}
}

type testContext struct {
	Name string
}

func TestContext(t *testing.T) {
	Ctx := testContext{Name: "test"}
	r := NewRouter()
	r.With(Ctx)

	r.GET("/", func(writer http.ResponseWriter, request *http.Request) {
		context := Context(request).(testContext)
		if context.Name != "test" {
			t.Errorf("Context was not passed to the handler")
		}
	})
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
}

func TestVars(t *testing.T) {
	r := NewRouter()
	r.GET("/user/{user}/{id}", func(writer http.ResponseWriter, request *http.Request) {
		vars := Vars(request)
		if user, ok := vars["user"]; !ok || user != "test" {
			t.Errorf("{user} was not passed to handler: want (test)")
		}
		if id, ok := vars["id"]; !ok || id != "10" {
			t.Errorf("{id} was not passed to handler: want (10)")
		}
	}).Where("user", "[a-z]+").Where("id", "[0-9]+")
	req := httptest.NewRequest("GET", "/user/test/10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
}
