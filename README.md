# gopherhole

# Usage

```
  config := gopherhole.NewConfiguration{}
  config.Port = '...'

  // ... or ...
  config := gopherhole.NewConfigurationFromFile("/path/to/myconfig.json")

  server := gopherhole.NewServer(config)
  server.Run()
```

# Configuration

## Settings and Defaults

| Attribute | Default Value | Description |
|-----------|---------------|-------------|
| RootDirectory | "/var/gopherhole" | The root path of your gopherhole. |
| Host | "localhost" | The hostname to use when building a directory file list response.  Clients will be instructed to make requests back to this domain, so you will want to change it from the default. |
| Port | 70 | The port to listen on, and to use when building a directory file list response. |
| MaxConnections | 0 | The maximum number of concurrent connections to allow.  The default, 0, means no limit.  It is recommended that you change this to a suitable non-zero value. |
| MapFileName | "gophermap" | The name of the file that, if present in a directory, will be served instead of a directory file list. |
| LogDisabled | false | If true, logging output will be supressed. |
| LogFile | n/a | If present, logs will be output to this file instead of STDOUT. |
| IdleTimeout | 60 | The number of seconds after which the connection will be closed if no query is made by the client.  0 means no timeout. |
| DefaultMimeType | "text/plain" | When a file's mimetype cannot be detected, it will be considered as this mime type. |
| MimeTypeIgnoreList | [] | An array of strings that represent mime types that should be ignored when building a directory list.  Entries can include the full type and subgroup ("text/plain"), or the type without a subgroup to exclude all mime-types in the group ("text/", note the trailing slash is required). |

## Configuration File

A configuration file must be in JSON, and may have any of the
attributes listed above.  Omitted attributes will use the default
values.

```
{
	"RootDirectory": "/var/gopherhole",
	"Host": "joshcom.net",
	"Port": 70,
	"MaxConnections": 5,
        "MapFileName": "gophermap",
        "LogDisabled": false,
	"LogFile": "/var/logs/gopherhole.log",
	"IdleTimeout": 30,
        "DefaultMineType": "text/plain",
	"MimeTypeIgnoreList": ["application/", "text/html"]
}
```

## Mime Types

The supported mime types will vary by system.  The [mime](https://golang.org/pkg/mime/) package is used to determine a file's mime type of its extension.  Review the documentation [here](https://golang.org/pkg/mime/#TypeByExtension) to see how mime types will be determined based on your OS.

# Gopher Maps

If a gopher map (specified by the `MapFileName` setting, defaulting to `gophermap`) is present in a directory.  This file is lightly processed, so that any line not detected have a <TAB> character will be presumed to be an inline-text line, and presented in the payload as such.

In other words, a gophermap with the contents (where <TAB> represents a tab character):

```
Welcome to my gopherhole!
0About this gopherhole<TAB>about.txt
1Phlog<TAB>phlog
```

...will be delivered (roughly) as:

```
iWelcome to my gopherhole!<TAB>(NOTHING)<TAB>nohost<TAB>0
0About this gopherhole<TAB>about.txt
1Phlog<TAB>phlog
```

The takeaway being:
* You do not need to format your inline-text lines as inline-text entities.  This will be done for you.
* Do not use the <TAB> character on inline-text lines (unless you've opted to format it yourself), or the previous bullet will not apply.

# Limitations

Currently, when files are served, they are loaded fully into memory, and then sent to the client.  They are not streamed.  (I intend to change this in the future.)
