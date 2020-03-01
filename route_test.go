package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var routesTest = []struct {
	p      string
	route  *Route
	expect bool
}{
	{"/user/testuser", NewRouter().GET("/user/{user}", testHandler).Where("user", "[a-z]+$"), true},
	{"/user/usettest10", NewRouter().GET("/user/{user}", testHandler).Where("user", "[a-z0-9]+$"), true},
	{"/user/usettest10", NewRouter().GET("/user/{user}", testHandler).Where("user", "[a-z]+$"), false},
}

func testHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Test handler"))
}

func TestRoute_Match(t *testing.T) {

	for _, r := range routesTest {
		req := httptest.NewRequest("GET", r.p, nil)
		if r.route.Match(req) != r.expect {
			t.Errorf("Error: Routes not match want %s got %s", r.route.where["user"], req.URL.Path)
		}
	}
}

func TestRoute_Where(t *testing.T) {
	r := NewRoute("/", "GET", testHandler)
	r.Where("user", ".*")
	if len(r.where) != 1 {
		t.Errorf("where was not added to route")
	}
}

func TestRoute_ExtractVars(t *testing.T) {
	r := NewRouter().GET("/user/{username}/{id}", testHandler).
		Where("username", ".*").
		Where("id", "[0-9]+")
	req := httptest.NewRequest("GET", "http://localhost:3000/user/test/10", nil)
	vars := r.extractVars(req)
	if username, ok := vars["username"]; !ok || username != "test" {
		t.Errorf("Error: fucked")
	}
}
