package gorilla_test

import (
	"testing"

	"github.com/gorilla/sessions"
	"ems/session/gorilla"
	"ems/session/test"
)

func TestAll(t *testing.T) {
	engine := sessions.NewCookieStore([]byte("something-very-secret"))
	manager := gorilla.New("_session", engine)
	test.TestAll(manager, t)
}
