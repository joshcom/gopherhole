package gopherhole

import (
	"testing"
)

func TestSelector_newSelector(t *testing.T) {
	t.Run("normal selector", func(t *testing.T) {
		selector := newSelector("/var/gopherhole", "phlog/")

		if selector.rootPath != "/var/gopherhole" {
			t.Errorf("Unexpected root path %s", selector.rootPath)
		}

		if selector.query != "phlog/" {
			t.Errorf("Unexpected query %s", selector.query)
		}
	})

	t.Run("strip trailing and leading spaces", func(t *testing.T) {
		selector := newSelector("/var/gopherhole", "/    \r\n")

		if selector.query != "/" {
			t.Errorf("Unexpected query %s", selector.query)
		}
	})
}

func TestSelector_path(t *testing.T) {
	t.Run("success paths", func(t *testing.T) {
		cases := []struct {
			root     string
			path     string
			expected string
		}{
			{
				"/var/gopherhole",
				"phlog/",
				"/var/gopherhole/phlog",
			},
			{
				"/",
				"phlog/",
				"/phlog",
			},
			{
				"/var/gopherhole",
				"/",
				"/var/gopherhole",
			},
			{
				"/var/gopherhole/",
				"//phlog//",
				"/var/gopherhole/phlog",
			},
			{
				"/var/gopherhole",
				"../gopherhole/phlog/",
				"/var/gopherhole/phlog",
			},
			{
				"/var/gopherhole",
				"../gopherhole/somefolder/../phlog/..",
				"/var/gopherhole",
			},
			{
				"/var/gopherhole",
				"",
				"/var/gopherhole",
			},
			{
				"/var/gopherhole",
				"~joshcom",
				"/var/gopherhole/~joshcom",
			},
			{
				"/var/gopherhole",
				"art/laptop.txt",
				"/var/gopherhole/art/laptop.txt",
			},
		}

		for _, testCase := range cases {
			selector := newSelector(testCase.root, testCase.path)
			fullPath, err := selector.path()

			if err != nil {
				t.Fatalf("Unexpected error %v with root=%s and path=%s",
					err,
					testCase.root,
					testCase.path)
			}

			if fullPath != testCase.expected {
				t.Errorf("Expected full path %s, but found %s",
					testCase.expected,
					fullPath)
			}
		}
	})

	t.Run("failure paths", func(t *testing.T) {
		cases := []struct {
			root string
			path string
		}{
			{
				"",
				"phlog/",
			},
			{
				"     ",
				"phlog/",
			},
			{
				"\n",
				"phlog/",
			},
			{
				"\r\n",
				"phlog/",
			},
			{
				" \r\n ",
				"phlog/",
			},
			{
				"/var/gopherhole",
				"../phlog",
			},
			{
				"/var/gopherhole",
				"..",
			},
			{
				"/var/gopherhole",
				"/..",
			},
			{
				"/var/gopherhole",
				"../",
			},
			{
				"/var/gopherhole",
				"/joshcom/phlog/../../../..",
			},
		}

		for _, testCase := range cases {
			selector := newSelector(testCase.root, testCase.path)
			fullPath, err := selector.path()

			if err == nil {
				t.Fatalf("Expected error with full path %s (root=%s and path=%s)",
					fullPath,
					testCase.root,
					testCase.path)
			}
		}
	})
}
