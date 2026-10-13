package store

import "testing"

func TestAddListDelete(t *testing.T) {
	s := New()
	before := len(s.List())

	e := s.Add("テストの波", MoodCalm)
	if e.ID == "" {
		t.Fatal("expected id")
	}
	if len(s.List()) != before+1 {
		t.Fatalf("expected %d entries, got %d", before+1, len(s.List()))
	}
	if s.List()[0].Body != "テストの波" {
		t.Fatalf("newest entry should be first")
	}

	if !s.Delete(e.ID) {
		t.Fatal("delete should succeed")
	}
	if s.Delete(e.ID) {
		t.Fatal("second delete should fail")
	}
	if len(s.List()) != before {
		t.Fatalf("expected %d entries after delete", before)
	}
}

func TestValidMood(t *testing.T) {
	if _, ok := ValidMood("calm"); !ok {
		t.Fatal("calm should be valid")
	}
	if _, ok := ValidMood("nope"); ok {
		t.Fatal("nope should be invalid")
	}
}
