package models

import (
	"encoding/json"
	"strings"
	"testing"
)

// Issue #271: ChatPermissions must be able to represent an explicit false on
// the wire. Plain bool + omitempty silently dropped false fields, which made
// callers think `CanSendMessages: false, CanSendOtherMessages: true` muted
// text while in fact only the true field reached Telegram.
func TestChatPermissions_ExplicitFalseSerialized(t *testing.T) {
	f, tr := false, true
	p := ChatPermissions{
		CanSendMessages:      &f,
		CanSendPhotos:        &tr,
		CanSendOtherMessages: &tr,
	}
	out, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	s := string(out)
	if !strings.Contains(s, `"can_send_messages":false`) {
		t.Fatalf("expected explicit can_send_messages:false, got %s", s)
	}
	if !strings.Contains(s, `"can_send_photos":true`) {
		t.Fatalf("expected can_send_photos:true, got %s", s)
	}
}

func TestChatPermissions_NilFieldOmitted(t *testing.T) {
	tr := true
	p := ChatPermissions{CanSendPhotos: &tr}
	out, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.Contains(string(out), "can_send_messages") {
		t.Fatalf("nil field should be omitted, got %s", out)
	}
}
