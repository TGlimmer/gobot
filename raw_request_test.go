package bot

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

type clientMock struct {
	requestURI string
	body       []byte
}

func (c *clientMock) Do(req *http.Request) (*http.Response, error) {
	c.requestURI = req.URL.RequestURI()
	if req.Body != nil {
		c.body, _ = io.ReadAll(req.Body)
	}
	resp := http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
	}
	return &resp, nil
}

type clientCaptureMock struct {
	method      string
	contentType string
	body        []byte
}

func (c *clientCaptureMock) Do(req *http.Request) (*http.Response, error) {
	c.method = req.Method
	c.contentType = req.Header.Get("Content-Type")
	if req.Body != nil && req.Body != http.NoBody {
		c.body, _ = io.ReadAll(req.Body)
	}
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
	}, nil
}

func Test_rawRequest_url(t *testing.T) {
	cm := &clientMock{}
	b := &Bot{
		token:  "XXX",
		client: cm,
	}

	err := b.rawRequest(context.Background(), "foo", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cm.requestURI != "/botXXX/foo" {
		t.Fatalf("unexpected requestURI: %s", cm.requestURI)
	}
}

func Test_rawRequest_url_testEnv(t *testing.T) {
	cm := &clientMock{}
	b := &Bot{
		token:           "XXX",
		client:          cm,
		testEnvironment: true,
	}

	err := b.rawRequest(context.Background(), "foo", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cm.requestURI != "/botXXX/test/foo" {
		t.Fatalf("unexpected requestURI: %s", cm.requestURI)
	}
}

// Issues #220, #224, #236: parameterless methods (GetMe, DeleteWebhook, ...)
// must not be wrapped in multipart. A zero-part multipart payload
// (\r\n--boundary--\r\n) is technically valid per RFC 2046 but the official
// Telegram backend and the self-hosted bot API server respond to it with an
// empty body, surfaced upstream as "unexpected end of JSON input". When
// params is nil the request is now sent as a bodyless POST with no
// Content-Type, which Telegram accepts.
func Test_rawRequest_nilParams_bodylessPost(t *testing.T) {
	cm := &clientCaptureMock{}
	b := &Bot{token: "XXX", client: cm}

	if err := b.rawRequest(context.Background(), "getMe", nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cm.body) != 0 {
		t.Fatalf("expected empty body, got %d bytes: %q", len(cm.body), cm.body)
	}
	if cm.contentType != "" {
		t.Fatalf("expected no Content-Type for parameterless request, got %q", cm.contentType)
	}
	if cm.method != http.MethodPost {
		t.Fatalf("expected POST, got %s", cm.method)
	}
}

// Same path as the above but with a typed-nil pointer (the shape callers hit
// via wrappers like DeleteWebhook(ctx, nil)). reflect.ValueOf on a typed nil
// pointer is not the same as on an untyped nil interface, so guard both.
func Test_rawRequest_typedNilParams_bodylessPost(t *testing.T) {
	cm := &clientCaptureMock{}
	b := &Bot{token: "XXX", client: cm}

	var params *DeleteWebhookParams // typed nil
	if err := b.rawRequest(context.Background(), "deleteWebhook", params, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cm.body) != 0 {
		t.Fatalf("expected empty body for typed-nil params, got %d bytes: %q", len(cm.body), cm.body)
	}
	if cm.contentType != "" {
		t.Fatalf("expected no Content-Type, got %q", cm.contentType)
	}
}

// With real params, the request still goes through the multipart path with
// the proper Content-Type and a non-empty body.
func Test_rawRequest_withParams_multipart(t *testing.T) {
	cm := &clientCaptureMock{}
	b := &Bot{token: "XXX", client: cm}

	params := &SendMessageParams{ChatID: int64(42), Text: "hello"}
	if err := b.rawRequest(context.Background(), "sendMessage", params, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.HasPrefix([]byte(cm.contentType), []byte("multipart/form-data; boundary=")) {
		t.Fatalf("expected multipart Content-Type, got %q", cm.contentType)
	}
	if !bytes.Contains(cm.body, []byte("hello")) {
		t.Fatalf("expected body to contain payload, got: %q", cm.body)
	}
}
