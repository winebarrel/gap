# gap

[![CI](https://github.com/winebarrel/gap/actions/workflows/ci.yml/badge.svg)](https://github.com/winebarrel/gap/actions/workflows/ci.yml)

Simple HTTP proxy with Google Oauth2 access token verification.

## Usage

```
gap --help
Usage: gap --backend=BACKEND --port=UINT --domain=STRING [flags]

Flags:
  -h, --help               Show help.
  -b, --backend=BACKEND    Backend URL to proxy ($GAP_BACKEND).
  -p, --port=UINT          Listening port ($GAP_PORT).
  -d, --domain=STRING      Allowed email domain ($GAP_DOMAIN).
      --version
```

```sh
$ go run ./cmd/gap -b https://example.comp -d winebarrel.jp -p 8080
```

```sh
# Get Oauth2 token and set it in an environment variable.
# e.g. https://developers.google.com/apps-script/reference/script/script-app?hl=ja#getOAuthToken()
$ export TOKEN='xxx'

$ curl  -H "x-gap-token: ${NOT_CORRECT_TOKEN}" localhost:8080
forbidden

~% curl -H "x-gap-token: ${TOKEN}" localhost:8080
<!doctype html>
<html>
<head>
    <title>Example Domain</title>
...
```
