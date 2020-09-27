package gopherhole

import (
	"testing"
)

func TestMime_guess(t *testing.T) {
	var m mimeType

	cases := []struct {
		path     string
		expected string
	}{
		{
			"/home/joshcom/textfile.txt",
			"text/plain",
		},
		{
			"textfile.txt",
			"text/plain",
		},
		{
			"barkingpuppy.gif",
			"image/gif",
		},
		{
			"mystory.phlog",
			"",
		},
	}

	for _, testCase := range cases {
		m.guess(testCase.path)
		if m != mimeType(testCase.expected) {
			t.Errorf("Expected '%s' to be mime type '%s', not '%s'.",
				testCase.path, testCase.expected, m)
		}
	}
}

func TestMime_string(t *testing.T) {
	cases := []struct {
		mime     mimeType
		expected string
	}{
		{
			mimeType("text/plain"),
			"text/plain",
		},
		{
			mimeType("image/gif"),
			"image/gif",
		},
		{
			mimeType(""),
			"",
		},
		{
			*new(mimeType),
			"",
		},
	}

	for _, testCase := range cases {
		name := testCase.mime.string()
		if name != testCase.expected {
			t.Errorf("Expected '%s', not '%s'.",
				name, testCase.expected)
		}
	}

}

func TestMime_isNone(t *testing.T) {
	m := mimeType("")
	if !m.isNone() {
		t.Error("Expected missing mime type, but one found.")
	}

	m = *new(mimeType)
	if !m.isNone() {
		t.Error("Expected missing mime type, but one found.")
	}

	m = mimeType("text/plain")
	if m.isNone() {
		t.Error("Expected mime type, but found none.")
	}
}

func TestMime_in(t *testing.T) {
	m := mimeType("text/html")
	t.Run("match conditions", func(t *testing.T) {
		conditions := [][]string{
			[]string{"text/"},
			[]string{"text/html"},
			[]string{"text/", "text/html"},
			[]string{"application/", "text/html"},
		}

		for _, condition := range conditions {
			if !m.in(condition) {
				t.Errorf("Mime '%s' expected to be in %v",
					m.string(),
					condition)
			}
		}
	})

	t.Run("non-match conditions", func(t *testing.T) {
		conditions := [][]string{
			[]string{"text"},
			[]string{"text/plain"},
			[]string{"text", "text/plain"},
			[]string{"application/", "text/plain"},
		}

		for _, condition := range conditions {
			if m.in(condition) {
				t.Errorf("Mime '%s' not expected to be in %v",
					m.string(),
					condition)
			}
		}
	})
}
