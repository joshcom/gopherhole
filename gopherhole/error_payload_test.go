package gopherhole

import (
	"strings"
	"testing"
)

func TestErrorPayload_build(t *testing.T) {
	payload := new(errorPayload)

	t.Run("build error", func(t *testing.T) {
		resource, _ := newQueryErrorResource("/fakepath")

		data, err := payload.build(&resource)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		dataStr := string(*data)
		formatedLine := "3Invalid query: /fakepath\t(NOTHING)\tnohost\t0\r\n"

		if strings.Index(dataStr, formatedLine) == -1 {
			t.Error("Unexpected selector format")
		}
		if !strings.HasSuffix(dataStr, string(payload.suffix())) {
			t.Error("Expected terminating suffix in payload.")
		}
	})
}
