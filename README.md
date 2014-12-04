# gotfavicon

Serves the favicon for a given URL.

A Go version of [getfavicon](https://github.com/potatolondon/getfavicon).

Favicons are cached on disk using a LRU cache.

## Installation

To run gotfavicon, you need to build the docker image:

```console
$ make
$ sudo docker run -d -p 8000:8000 -t basiclytics/gotfavicon
```

## Quickstaty

```html
<img src="http://gotfaviconfrontend/?url=http://www.google.com" alt="Google" width="16" height="16" />
```

## License

MIT
