package middleware

import (
	"crypto/sha256"
	"testing"
)

func TestAuthToken(t *testing.T) {
	raw := []byte(`{"a":1}`)
	key := "pi"
	h := sha256.New()

	token, md, err := verify(AuthToken(raw, key, h), key, h)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if string(raw) != token {
		t.Logf("%s expected but got %s", string(raw), token)
		t.FailNow()
	}
	if md["a"] != 1.0 {
		t.Logf("1 expected but got %f", md["a"])
		t.FailNow()
	}
}

func TestSplit(t *testing.T) {
	s := "hello.world"
	s1, s2, err := split(s)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if s1 != "hello" {
		t.Logf("hello expected but got %s", s1)
		t.FailNow()
	}
	if s2 != "world" {
		t.Logf("world expected but got %s", s2)
		t.FailNow()
	}
}
