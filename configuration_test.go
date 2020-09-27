package gopherhole

import (
	"testing"
)

func TestConfiguration_NewConfiguration(t *testing.T) {
	config := NewConfiguration()

	if ROOT_DIRECTORY != config.RootDirectory {
		t.Errorf("Expected directory '%s', found '%s'.",
			ROOT_DIRECTORY,
			config.RootDirectory)
	}

	if PORT != config.Port {
		t.Errorf("Expected port %d, found %d.", PORT, config.Port)
	}

	if MAX_CONNECTIONS != config.MaxConnections {
		t.Errorf("Expected max connections %d, found %d.",
			MAX_CONNECTIONS,
			config.MaxConnections)
	}
}

func TestConfiguration_NewConfigurationFromFile(t *testing.T) {
	t.Run("test loaded config", func(t *testing.T) {
		config, err := NewConfigurationFromFile("testdata/config/config.json")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if config.RootDirectory != "testdata/mygopherhole" {
			t.Error("Configuration not loaded as expected.")
		}

		if config.MapFileName != "gophermap" {
			t.Error("Expected missing config setting to have default.")
		}
	})

	t.Run("test missing file", func(t *testing.T) {
		_, err := NewConfigurationFromFile("testdata/config/config.xml")
		if err == nil {
			t.Error("Expected an error on missing file.")
		}
	})

	t.Run("test non-JSON file", func(t *testing.T) {
		_, err := NewConfigurationFromFile("testdata/mygopherhole/art/laptop.txt")
		if err == nil {
			t.Error("Expected an error non-JSON file.")
		}
	})
}
