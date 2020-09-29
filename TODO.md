# TODO

## Additional TODOs
0. Makefile
1. What happens with various invalid config settings? (selector and root path, mime type)
2. Decorate incomplete lines of gophermap with default hostname and port.
3. Review the period situation (https://tools.ietf.org/html/rfc1436)
4. Make sure all file absolute paths are in root (and make an absolute path out of the provided root config).
5. Coverage report and CI

## Future projects
* Streaming writes, instead of loading the file in memory.
* Additional mime-types file. (/etc/mime.types)
* Use channel, instead of sleep, when connections max out.
