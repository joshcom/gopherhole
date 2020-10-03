# TODO

## Additional TODOs
1. What happens with various invalid config settings? (selector and root path, mime type)
2. Review the period situation (https://tools.ietf.org/html/rfc1436)
3. Make sure all file absolute paths are in root (and make an absolute path out of the provided root config).
4. Coverage report and CI

## Future projects
* Streaming writes, instead of loading the file in memory.
* Additional mime-types file. (/etc/mime.types)
* Use channel, instead of sleep, when connections max out.
