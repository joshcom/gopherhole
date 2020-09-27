package gopherhole

import (
	"testing"
)

func TestEntity_EntityForMimeType(t *testing.T) {
	var m EntityType
	var ok bool

	t.Run("direct match", func(t *testing.T) {
		m, ok = EntityForMimeType("image/gif")
		if m != GifEntity {
			t.Errorf("Expected %s, found %s", string(GifEntity), string(m))
		}
	})

	t.Run("group match", func(t *testing.T) {
		m, ok = EntityForMimeType("text/html")
		if m != HtmlEntity {
			t.Errorf("Expected %s, found %s", string(TextFileEntity), string(m))
		}
	})

	t.Run("no matches", func(t *testing.T) {
		cases := []string{
			"text/",
			"text/   ",
			"text",
			"sports/golf",
			"sports",
			"html",
			"gif/images",
			"  gif/images",
			"gif/images   ",
		}

		for _, testCase := range cases {
			m, ok = EntityForMimeType(testCase)
			if ok == true {
				t.Errorf("Expected missing entity, but found %s for %s.", string(m), testCase)
			}
		}
	})
}
