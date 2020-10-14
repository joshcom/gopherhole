package gopherhole

import (
	"strings"
	"testing"
)

func TestMapfilePayload_build(t *testing.T) {
	payload := newMapfilePayload("joshcom.net", 70, "testdata/mygopherhole")

	t.Run("build response for file", func(t *testing.T) {
		path := "testdata/mygopherhole/phlog/gophermap"
		resource, err := newResource(path)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		reader, err := payload.build(&resource)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		dataStr := readAllString(*reader)

		if !strings.HasSuffix(dataStr, string(payload.suffix())) {
			t.Error("Terminating suffix should be applied to mapfiles.")
		}

		// Spot check a few lines.
		expectedRows := []string{
			"ijoshcom.net - PHLOG\t(NOTHING)\tnohost\t0\r\n",
			"i\t(NOTHING)\tnohost\t0\r\n",
			"1^^Top of the gopherhole^^\t/\tjoshcom.net\t70\r\n",
		}

		for _, expected := range expectedRows {
			if strings.Index(dataStr, expected) < 0 {
				t.Errorf("Expected line not found:\n%s", expected)
			}
		}
	})
}
