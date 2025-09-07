# gap

[![CI](https://github.com/winebarrel/gap/actions/workflows/ci.yml/badge.svg)](https://github.com/winebarrel/gap/actions/workflows/ci.yml)

Simple HTTP proxy with Google Oauth2 access token verification.

## Usage

```
Usage: gap --backend=BACKEND --port=UINT --header-name=STRING --allow-list=ALLOW-LIST,... [flags]

Flags:
  -h, --help                  Show help.
  -b, --backend=BACKEND       Backend URL ($GAP_BACKEND).
  -p, --port=UINT             Listening port ($GAP_PORT).
  -n, --header-name=STRING    Header name to pass the access token
                              ($GAP_HEADER).
  -e, --allow-list=ALLOW-LIST,...
                              Allowed email list that may contain wildcards
                              ($GAP_ALLOW_LIST).
      --version
```

```sh
$ go run ./cmd/gap -b https://example.com -e '*@example.com' -p 8080 -n x-my-gap-token

# When using Docker:
# docker run --rm -p 8080:8080 ghcr.io/winebarrel/gap -b https://example.com -e '*@example.com' -p 8080 -n x-my-gap-token
```

```sh
# Get Oauth2 token from Apps Script and set it in an environment variable.
# e.g. https://developers.google.com/apps-script/reference/script/script-app#getOAuthToken()
$ export TOKEN='xxx'

$ curl  -H "x-my-gap-token: ${NOT_CORRECT_TOKEN}" localhost:8080
forbidden

~% curl -H "x-my-gap-token: ${TOKEN}" localhost:8080
<!doctype html>
<html>
<head>
    <title>Example Domain</title>
...
```
