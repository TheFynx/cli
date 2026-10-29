package netcheck

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCheckVerdicts(t *testing.T) {
	ok := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ok.Close()
	bad := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer bad.Close()

	up, detail, err := Check(ok.URL, time.Second)
	if !up || err != nil || detail == "" {
		t.Fatalf("200 must be up: %v %q %v", up, detail, err)
	}
	up, _, err = Check(bad.URL, time.Second)
	if up || err != nil {
		t.Fatalf("500 must be explicitly down, network failure: %v %v", up, err)
	}
	_, _, err = Check("http://127.0.0.1:1/x", time.Second)
	if err == nil {
		t.Fatal("connection failure must surface as an error")
	}
}
