package rrouter

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestUse(t *testing.T) {
	s := New()
	mw := func(fn http.HandlerFunc) http.HandlerFunc {
		return fn
	}
	c := len(s.middlewares)

	s.Use(mw)

	if len(s.middlewares) != c+1 {
		t.Error("expected Use() to increase the number of items in the stack")
	}
}

func TestWrap(t *testing.T) {
	s := New()

	var middlewareCalled bool
	mw := func(fn http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			middlewareCalled = true
			fn(w, r)
		}
	}
	s.Use(mw)

	var handlerCalled bool
	hn := func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	}

	wrapped := s.Wrap(hn)
	req := httptest.NewRequest("GET", "/example", nil)
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(plainHandler(wrapped))
	handler.ServeHTTP(w, req)

	if !handlerCalled {
		t.Error("expected handler to have been called")
	}

	if !middlewareCalled {
		t.Error("expected middleware to have been called")
	}
}

func TestWrap_Ordering(t *testing.T) {
	s := New()

	var firstCallAt *time.Time
	var secondCallAt *time.Time
	var thirdCallAt *time.Time
	var fourthCallAt *time.Time
	var handlerCallAt *time.Time

	first := func(fn http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ts := time.Now()
			firstCallAt = &ts
			fn(w, r)
		}
	}

	second := func(fn http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ts := time.Now()
			secondCallAt = &ts
			fn(w, r)
		}
	}
	third := func(fn http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ts := time.Now()
			thirdCallAt = &ts
			fn(w, r)
		}
	}
	fourth := func(fn http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ts := time.Now()
			fourthCallAt = &ts
			fn(w, r)
		}
	}

	s.Use(first)
	s.Use(second)
	s.Use(third)
	s.Use(fourth)

	hn := func(w http.ResponseWriter, r *http.Request) {
		ts := time.Now()
		handlerCallAt = &ts
	}

	wrapped := s.Wrap(hn)
	req := httptest.NewRequest("GET", "/example", nil)
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(plainHandler(wrapped))
	handler.ServeHTTP(w, req)

	if firstCallAt == nil || secondCallAt == nil || thirdCallAt == nil || fourthCallAt == nil || handlerCallAt == nil {
		t.Fatal("failed to call one or more functions")
	}

	if firstCallAt.After(*secondCallAt) || firstCallAt.After(*thirdCallAt) || firstCallAt.After(*fourthCallAt) || firstCallAt.After(*handlerCallAt) {
		t.Error("failed to call first middleware first")
	}

	if fourthCallAt.Before(*thirdCallAt) || fourthCallAt.Before(*secondCallAt) || fourthCallAt.After(*handlerCallAt) {
		t.Error("failed to call fourth middleware last before the handler")
	}

	if secondCallAt.After(*thirdCallAt) {
		t.Error("expected second middleware to come before the third")
	}
}

func TestWrap_WhenEmpty(t *testing.T) {
	s := New()
	hn := func(w http.ResponseWriter, r *http.Request) {}
	w := s.Wrap(hn)

	if reflect.ValueOf(hn).Pointer() != reflect.ValueOf(w).Pointer() {
		t.Error("expected that Wrap() would return the given function when stack is empty")
	}
}

func plainHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)
	}
}

func TestSubroute(t *testing.T) {

	sr := New()
	var handlerCalled bool
	hn := func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	}
	sr.GET("/subroute", hn)
	req := httptest.NewRequest("GET", "/parent/heheh/subroute", nil)
	w := httptest.NewRecorder()

	router := New()
	router.Subroute("/parent/:parent", sr)
	router.ServeHTTP(w, req)

	if !handlerCalled {
		t.Error("expected handler to have been called")
	}

}
