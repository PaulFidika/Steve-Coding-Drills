This is an example of native ES modules, as opposed to bundled sites.

Note that index.html cannot be opened from the file system, but needs a webserver, in order to bypass CORS protection in the browser. RUn `python -m http.server` to start a simple webserver.

Circular-imports are allowed as well.
