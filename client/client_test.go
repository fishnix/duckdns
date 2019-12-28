package client

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {
	testClient := NewClient("foodom", "111111-222222-333333-444444", true)
	cType := reflect.TypeOf(testClient).String()
	if cType != "*client.Client" {
		t.Errorf("expected client to be '*client.Client', got '%s'", cType)
	}
}

func TestSetHTTPClient(t *testing.T) {
	testClient := NewClient("foodom", "111111-222222-333333-444444", true)
	testClient.SetHTTPClient(nil)
	if testClient.HTTPClient != nil {
		t.Errorf("expected http client to be 'nil', got '%+v'", testClient.HTTPClient)
	}
}

func TestSetDomain(t *testing.T) {
	testClient := NewClient("foodom", "111111-222222-333333-444444", true)
	testClient.SetDomain("")
	if testClient.Domain != "" {
		t.Errorf("expected domain to be '', got '%s'", testClient.Domain)
	}
}

func TestSetToken(t *testing.T) {
	testClient := NewClient("foodom", "111111-222222-333333-444444", true)
	testClient.SetToken("")
	if testClient.Token != "" {
		t.Errorf("expected token to be '', got '%s'", testClient.Token)
	}
}

func TestSetBaseURL(t *testing.T) {
	testClient := NewClient("foodom", "111111-222222-333333-444444", true)
	testClient.SetBaseURL("")
	if testClient.BaseURL != "" {
		t.Errorf("expected base URL to be '', got '%s'", testClient.BaseURL)
	}
}

func TestSetVerbose(t *testing.T) {
	testClient := NewClient("foodom", "111111-222222-333333-444444", true)
	testClient.SetVerbose(true)
	if testClient.Verbose != true {
		t.Errorf("expected verbose to be 'true', got '%t'", testClient.Verbose)
	}
}

func TestUpdate(t *testing.T) {

	type test struct {
		domain string
		token  string
		status int
		body   []byte
		err    error
	}

	tests := []test{
		test{
			domain: "foo",
			token:  "goodtoken",
			status: http.StatusOK,
			body:   []byte("OK"),
			err:    nil,
		},
		test{
			domain: "bar",
			token:  "badtoken",
			status: http.StatusOK,
			body:   []byte("KO"),
			err:    errors.New("got failure response to update: KO"),
		},
		test{
			domain: "baz",
			token:  "goodtoken",
			status: http.StatusBadRequest,
			body:   []byte("Bad Request"),
			err:    errors.New("unexpected response from update: Bad Request"),
		},
	}

	for _, set := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Logf("got request: %+v", r)
			w.WriteHeader(set.status)
			w.Write(set.body)
		}))
		defer ts.Close()

		testClient := NewClient(set.domain, set.token, false)
		testClient.SetBaseURL(ts.URL + "/update?domains=%s&token=%s")
		err := testClient.Update("", "")
		if err == nil && set.err == nil {
			t.Log("got expected nil error response")
		} else if err == nil && set.err != nil {
			t.Errorf("expected '%s' error, got '%s'", set.err, err)
		} else if err != nil && set.err == nil {
			t.Errorf("expected '%s' error, got '%s'", set.err, err)
		} else if err.Error() != set.err.Error() {
			t.Errorf("expected '%s' error, got '%s'", set.err, err)
		}

		testClient.SetVerbose(true)
		err = testClient.Update("", "")
		if err == nil && set.err == nil {
			t.Log("got expected nil error response")
		} else if err == nil && set.err != nil {
			t.Errorf("expected '%s' error, got '%s'", set.err, err)
		} else if err != nil && set.err == nil {
			t.Errorf("expected '%s' error, got '%s'", set.err, err)
		} else if err.Error() != set.err.Error() {
			t.Errorf("expected '%s' error, got '%s'", set.err, err)
		}

		err = testClient.Update("1.2.3.4", "")
		if err == nil && set.err == nil {
			t.Log("got expected nil error response")
		} else if err == nil && set.err != nil {
			t.Errorf("expected '%s' error, got '%s'", set.err, err)
		} else if err != nil && set.err == nil {
			t.Errorf("expected '%s' error, got '%s'", set.err, err)
		} else if err.Error() != set.err.Error() {
			t.Errorf("expected '%s' error, got '%s'", set.err, err)
		}

		err = testClient.Update("", "1111:2222:3333:4444:5555:6666:aaaa:bbbb")
		if err == nil && set.err == nil {
			t.Log("got expected nil error response")
		} else if err == nil && set.err != nil {
			t.Errorf("expected '%s' error, got '%s'", set.err, err)
		} else if err != nil && set.err == nil {
			t.Errorf("expected '%s' error, got '%s'", set.err, err)
		} else if err.Error() != set.err.Error() {
			t.Errorf("expected '%s' error, got '%s'", set.err, err)
		}
	}
}
