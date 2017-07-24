package ring

import (
	"io"
	"testing"
)

func TestWriter(t *testing.T) {
	b := make([]byte, 5)
	steps := []struct {
		p string // input
		s string // state of b
	}{
		{"abc", "abc\x00\x00"},
		{"de", "abcde"},
		{"fghijk", "kghij"},
		{"aaaaa", "aaaaa"},
	}

	w := NewWriter(b, 0)
	for i, step := range steps {
		n, err := w.Write([]byte(step.p))
		if err != nil {
			// should newer happen
			t.Fatal(err)
		}
		if m := len(step.p); n != m {
			t.Errorf("#%d: wrote %d bytes, want %d bytes", i, n, m)
		}
		if s := string(b); s != step.s {
			t.Errorf("#%d: got %q, want %q", i, s, step.s)
		}
	}
}

func TestLimitWriter(t *testing.T) {
	b := make([]byte, 5)
	w := NewLimitWriter(b, 0, 4)
	n, err := w.Write([]byte("abcde"))
	if n != 4 {
		t.Errorf("wrote %d bytes, want 4 bytes", n)
	}
	if err != io.ErrShortWrite {
		t.Errorf("got error %v, want error %v", err, io.ErrShortWrite)
	}
	want := "abcd\x00"
	if s := string(b); s != want {
		t.Errorf("got %q, want %q", s, want)
	}
}

func TestReader(t *testing.T) {
	b := []byte("abcde")
	steps := []struct {
		n int    // bytes to read
		s string // expected output
	}{
		{1, "a"},
		{4, "bcde"},
		{5, "abcde"},
		{6, "abcdea"},
	}

	r := NewReader(b, 0)
	for i, step := range steps {
		p := make([]byte, step.n)
		n, err := r.Read(p)
		if err != nil {
			// should newer happen
			t.Fatal(err)
		}
		if n != step.n {
			t.Errorf("#%d: read %d bytes, want %d bytes", i, n, step.n)
		}
		if s := string(p); s != step.s {
			t.Errorf("#%d: got %q, want %q", i, s, step.s)
		}
	}
}

func TestLimitReader(t *testing.T) {
	b := []byte("abcde")
	p := make([]byte, 5)
	r := NewLimitReader(b, 0, 4)
	n, err := r.Read(p)
	if n != 4 {
		t.Errorf("read %d bytes, want 4 bytes", n)
	}
	if err != io.EOF {
		t.Errorf("got error %v, want error %v", err, io.EOF)
	}
	want := "abcd\x00"
	if s := string(p); s != want {
		t.Errorf("got %q, want %q", s, want)
	}
}
