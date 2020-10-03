# gopher

# Usage

```
  config := gopher.NewConfiguration{}
  config.Port = '...'

  // ... or ...
  config := gopher.NewConfigurationFromFile("/path/to/myconfig.json")

  server := gopher.NewServer(config)
  server.Run()
```

# Configuration

See the parent README for details on configuration settings.

