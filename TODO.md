# TODO

## Additional TODOs
1. What happens with various invalid config settings? (selector and root path, mime type)
2. Server: Better tests, better overall abstraction?
3. Decorate incomplete lines of gophermap with default hostname and port.
4. Review the period situation (https://tools.ietf.org/html/rfc1436)
5. How do we handle file permissions?
6. Make sure all file absolute paths are in root (and make an absolute path out of the provided root config).
7. Do not directly serve files forbidden by map restrictions.

## Deployment
* How to package this?
* README
* How to run as daemon, and log to syslog?
* Run as a specific user?
* PID file?
* Gracefully wait for handlers to complete on shutdown.

## Future projects
* Streaming writes, instead of loading the file in memory.
* Option to require gophermaps
* Additional mime-types file. (/etc/mime.types)
