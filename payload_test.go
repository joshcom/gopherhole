package gopherhole

import (
	"strings"
	"testing"
)

func TestPayload_pack(t *testing.T) {
	data := []byte("This is a response.")
	payload := payloadImpl{}
	packedPayload := payload.pack(&data)

	dataExpected := []byte("This is a response.")
	packedExpected := []byte("This is a response.\r\n.\r\n")

	if string(data) != string(dataExpected) {
		t.Error("Original data mutated.")
	}

	if string(*packedPayload) != string(packedExpected) {
		t.Error("Data not packed as expected.")
	}
}

func TestPayload_isEntityRow(t *testing.T) {
	payload := payloadImpl{}
	t.Run("detects formatted entity rows", func(t *testing.T) {
		goodRows := [][]byte{
			[]byte("0My Computer\tcomputer.txt\tjoshcom.et\t70\r\n"),
			[]byte("iThis is great!\t(NOTHING)\tnohost\t0\r\n"),
			[]byte("1Phlog\tphlog/\tjoshcom.net\t70\r\n"),
			[]byte("1i\t \twell\t1\r\n"),
		}

		for _, row := range goodRows {
			result := payload.isEntityRow(&row)
			if !result {
				t.Errorf("Row expected to be an entity:\n%s", row)
			}
		}
	})

	t.Run("detects non-entity rows", func(t *testing.T) {
		badRows := [][]byte{
			[]byte("0My Computer\r\n"),
			[]byte("iThis is great!"),
			[]byte("Phlog phlog/ joshcom.net 70\r\n"),
			[]byte("3This"),
			[]byte("i        nohost    0\r\n"),
			[]byte("i         nohost    0\r\n"),
		}

		for _, row := range badRows {
			result := payload.isEntityRow(&row)
			if result {
				t.Errorf("Row not expected to be an entity:\n%s", row)
			}
		}
	})
}

func TestPayload_buildInlineTextEntityRow(t *testing.T) {
	payload := payloadImpl{}
	row := payload.buildInlineTextEntityRow([]byte("This is the row\r\n"))
	expected := "iThis is the row\t(NOTHING)\tnohost\t0\r\n"

	if string(*row) != expected {
		t.Errorf("Unexpcted row format: \n%s", row)
	}
}

func TestPayload_buildResourceEntityRow(t *testing.T) {
	payload := payloadImpl{
		host:        "joshcom.net",
		port:        70,
		defaultMime: "text/plain",
	}

	t.Run("directory resource entity", func(t *testing.T) {
		path := "testdata/mygopherhole/art"
		res, err := newResource(path)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}
		row := payload.buildResourceEntityRow(&res)
		expected := "1art\ttestdata/mygopherhole/art\tjoshcom.net\t70\r\n"
		if string(*row) != expected {
			t.Errorf("Unexpcted row format: \n%s", row)
		}
	})

	t.Run("file resource entity", func(t *testing.T) {
		path := "testdata/mygopherhole/art/laptop.txt"
		res, err := newResource(path)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}
		row := payload.buildResourceEntityRow(&res)
		expected := "0laptop.txt\ttestdata/mygopherhole/art/laptop.txt\tjoshcom.net\t70\r\n"
		if string(*row) != expected {
			t.Errorf("Unexpcted row format: \n%s", row)
		}
	})

	t.Run("file resource entity without detectable mime type", func(t *testing.T) {
		path := "testdata/mygopherhole/phlog/gophermap"
		res, err := newResource(path)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}
		row := payload.buildResourceEntityRow(&res)
		expected := "0gophermap"
		if strings.Index(string(*row), expected) != 0 {
			t.Errorf("Unexpcted row format: \n%s", row)
		}
	})

	t.Run("error entity", func(t *testing.T) {
		path := "testdata/mygopherhole/art/laptop2.txt"
		res, _ := newResource(path)
		row := payload.buildResourceEntityRow(&res)
		expected := "3Requested resource not found.\t(NOTHING)\tnohost\t0\r\n"
		if string(*row) != expected {
			t.Errorf("Unexpcted row format: \n%s", row)
		}
	})
}

func TestPayload_buildErrorEntityRow(t *testing.T) {
	payload := payloadImpl{}
	row := payload.buildErrorEntityRow("File not found.")
	expected := "3File not found.\t(NOTHING)\tnohost\t0\r\n"

	if string(*row) != expected {
		t.Errorf("Unexpcted row format: \n%s", row)
	}
}

func TestPayload_buildEntityRow(t *testing.T) {
	payload := payloadImpl{}
	row := payload.buildEntityRow(DirectoryEntity, "Phlog", "phlog/", "joshcom.net", 70)
	expected := "1Phlog\tphlog/\tjoshcom.net\t70\r\n"

	if string(*row) != expected {
		t.Errorf("Unexpcted row format: \n%s", row)
	}
}
