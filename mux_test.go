package pat

import (
	"github.com/bmizerany/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPatMatch(t *testing.T) {
	params, splat, ok := (&patHandler{"/foo/:name", nil}).try("/foo/bar")
	assert.Equal(t, true, ok)
	assert.Equal(t, Params{":name": "bar"}, params)
	assert.Equal(t, splat, "")

	params, splat, ok = (&patHandler{"/foo/:name/baz", nil}).try("/foo/bar")
	assert.Equal(t, false, ok)
	assert.Equal(t, splat, "")

	params, splat, ok = (&patHandler{"/foo/:name/baz", nil}).try("/foo/bar/baz")
	assert.Equal(t, true, ok)
	assert.Equal(t, Params{":name": "bar"}, params)
	assert.Equal(t, splat, "")

	params, splat, ok = (&patHandler{"/foo/:name/baz/:id", nil}).try("/foo/bar/baz")
	assert.Equal(t, false, ok)
	assert.Equal(t, splat, "")

	params, splat, ok = (&patHandler{"/foo/:name/baz/:id", nil}).try("/foo/bar/baz/123")
	assert.Equal(t, true, ok)
	assert.Equal(t, Params{":name": "bar", ":id": "123"}, params)
	assert.Equal(t, splat, "")

	params, splat, ok = (&patHandler{"/foo/:name/baz/:name", nil}).try("/foo/bar/baz/123")
	assert.Equal(t, true, ok)
	assert.Equal(t, Params{":name": "123"}, params)
	assert.Equal(t, splat, "")

	params, splat, ok = (&patHandler{"/foo/::name", nil}).try("/foo/bar")
	assert.Equal(t, true, ok)
	assert.Equal(t, Params{"::name": "bar"}, params)
	assert.Equal(t, splat, "")

	params, splat, ok = (&patHandler{"/foo/x:name", nil}).try("/foo/bar")
	assert.Equal(t, false, ok)
	assert.Equal(t, splat, "")

	params, splat, ok = (&patHandler{"/foo/x:name", nil}).try("/foo/xbar")
	assert.Equal(t, true, ok)
	assert.Equal(t, Params{":name": "bar"}, params)
	assert.Equal(t, splat, "")

	params, splat, ok = (&patHandler{"/foo/", nil}).try("/foo/bar/baz")
	assert.Equal(t, true, ok)
	assert.Equal(t, splat, "bar/baz")

	params, splat, ok = (&patHandler{"/foo/", nil}).try("/foo/bar")
	assert.Equal(t, true, ok)
	assert.Equal(t, splat, "bar")
}

func TestPatRoutingHit(t *testing.T) {
	p := New()

	var ok bool
	p.Get("/foo/:name", HandlerFunc(func(p Params, _ string) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ok = true
			t.Logf("%#v", r.URL.Query())
			assert.Equal(t, "keith", p[":name"])
		})
	}))

	r, err := http.NewRequest("GET", "/foo/keith?a=b", nil)
	if err != nil {
		t.Fatal(err)
	}

	p.ServeHTTP(nil, r)

	assert.T(t, ok)
}

func TestPatRoutingNoHit(t *testing.T) {
	p := New()

	var ok bool
	p.Post("/foo/:name", HandlerFlat(func(p Params, splat string, w http.ResponseWriter, r *http.Request) {
		ok = true
	}))

	r, err := http.NewRequest("GET", "/foo/keith", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	p.ServeHTTP(rr, r)

	assert.T(t, !ok)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}
