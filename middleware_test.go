package middleware

import (
	"net/http"
	"testing"
)

var execOrder []string
var testCascade *testing.T

func middleA(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			execOrder = append(execOrder, "one")
			next.ServeHTTP(w, r)
			if len(execOrder) != 4 {
				testCascade.Errorf("Middlewares did not execute recursively (mA)")
			}
		},
	)
}

func middleB(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			execOrder = append(execOrder, "two")
			next.ServeHTTP(w, r)
			if len(execOrder) != 4 {
				testCascade.Errorf("Middlewares did not execute recursively (mB)")
			}
		},
	)
}

func middleC(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			execOrder = append(execOrder, "three")
			next.ServeHTTP(w, r)
			if len(execOrder) != 4 {
				testCascade.Errorf("Middlewares did not execute recursively (mC)")
			}
		},
	)
}

// TestCascase tests the Cascade function for correct order, and completeness of execution
func TestCascade(t *testing.T) {
	testCascade = t

	handler := Cascade(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				execOrder = append(execOrder, "four")
			},
		),
		middleA,
		middleB,
		middleC,
	)

	handler.ServeHTTP(nil, nil)
	if len(execOrder) != 4 {
		t.Errorf("Incorrect number of middlewares executed")
	}
	if execOrder[0] != "one" || execOrder[1] != "two" || execOrder[2] != "three" || execOrder[3] != "four" {
		t.Errorf("Middlewares executed in an unexpected order")
	}
}
