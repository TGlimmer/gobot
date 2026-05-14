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
// were shipping a multipart Content-Type with an empty body, which both the
// official Telegram backend and self-hosted bot API server reject with an
// empty response (surfaced as "unexpected end of JSON input"). The multipart
// writer must always emit a closing boundary, even when params is nil.
func Test_rawRequest_nilParams_writesClosingBoundary(t *testing.T) {
	cm := &clientMock{}
	b := &Bot{token: "XXX", client: cm}

	if err := b.rawRequest(context.Background(), "getMe", nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cm.body) == 0 {
		t.Fatalf("expected non-empty body with closing multipart boundary, got empty body")
	}
	// A finalized multipart body always ends with "--<boundary>--\r\n".
	if !bytes.HasSuffix(cm.body, []byte("--\r\n")) {
		t.Fatalf("expected body to end with closing multipart boundary, got:\n%q", cm.body)
	}
}

// Same path as the above but with a typed-nil pointer (the shape callers hit
// via wrappers like DeleteWebhook(ctx, nil)). reflect.ValueOf on a typed nil
// pointer is not the same as on an untyped nil interface, so guard both.
func Test_rawRequest_typedNilParams_writesClosingBoundary(t *testing.T) {
	cm := &clientMock{}
	b := &Bot{token: "XXX", client: cm}

	var params *DeleteWebhookParams // typed nil
	if err := b.rawRequest(context.Background(), "deleteWebhook", params, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.HasSuffix(cm.body, []byte("--\r\n")) {
		t.Fatalf("expected body to end with closing multipart boundary, got:\n%q", cm.body)
	}
}
