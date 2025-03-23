# `temple`

So I needed template executor for nginx configs that would go into production image.

I ended up writing this tool. It has zero dependencies. Output binary is around 2 MB.

Why not ...

- [`envsubst`](https://man7.org/linux/man-pages/man1/envsubst.1.html) - does simple substitution and does not support default values, conditions, or loops
- [`gomplate`](https://github.com/hairyhenderson/gomplate) - because it's too fat! Seriously, guys, 100 megs for a template renderer?!
- `sed`/`awk` - too much fuss

Tested with Go 1.10.7, 1.13, 1.16-1.19, 1.23.6

## Installation

### Pre-built binary

```sh
curl -L https://github.com/vegas503/temple/releases/latest/download/temple \
    -o /usr/bin/temple && \
    chmod +x /usr/bin/temple
```

### Install using go

```sh
# Build with default flags (~3.8 MB):
go install github.com/vegas503/temple/cmd/temple@latest

# Build with flags to produce a smaller binary (~2.3 MB):
go install \
    -trimpath \
    -ldflags '-s -w' \
    -gcflags=all='-B -l -wb=false' \
    github.com/vegas503/temple/cmd/temple@latest
```

## Usage

All the heavy lifting is done by Go's [text/template](https://pkg.go.dev/text/template) package. Go see their docs first!

In addition to the built-ins, `temple` provides some additional functions for common scenarios.

Some functions (e. g. `contains`, `split`', `join`, and `replace`) have their arguments swapped.
This is done to make it easier to use [function pipelines](https://pkg.go.dev/text/template#hdr-Examples).

For instance, `{{ $user := (split "," (env "USERS")) }}` can be written as `{{ $user := env "USERS" | split "," }}`.

### Functions

- `env STRING` - gets env var value, throws error if env var is not set

    ```sh
    $ echo '{{ env "USER" }}' | USER=skaarj temple
    skaarj
    ```

- `envdefault STRING STRING` - gets env var value, returns default value if empty or unset

    ```sh
    $ echo '{{ envdefault "SOME_UNSET_VAR" "foo" }}' | temple
    foo
    ```

- `split DELIM STRING` - splits a string with a delimiter, trims whitespaces and returns array with all empty elements removed

    ```sh
    $ echo '{{ range $v := (split "," (env "ITEMS")) }}{{ $v }}{{ end }}' | ITEMS="A, B ,  C " temple
    ABC
    $ echo '{{ range $v := env "ITEMS" | split "," }}{{ $v }}{{ end }}' | ITEMS="A, B ,  C " temple
    ABC
    ```

- `contains SUBSTRING STRING` - proxy for [strings.Contains](https://pkg.go.dev/strings#Contains) with arguments swapped

- `join GLUE ARRAY` - proxy for [strings.Join](https://pkg.go.dev/strings#Join) with arguments swapped

- `coalesce STRING...` - returns first non-empty string

    ```sh
    $ echo '{{ coalesce "" "two" "three" }}' | temple
    two
    ```

- `chain ARRAY...` - concatenates arrays

    ```sh
    $ cat temple.tpl
    {{- $colors := env "COLORS" | split "," }}
    {{- $numbers := env "NUMBERS" | split "," }}
    {{- range $item := (chain $colors $numbers) }}
      * {{ $item }}
    {{- end }}
    $ COLORS=red,green,blue NUMBERS=34,42 temple -i temple.tpl
      * red
      * green
      * blue
      * 34
      * 42

- `uniq ARRAY` - returns array with all duplicate elements removed

    ```sh
    $ echo '{{ range $v := (uniq (split "," (env "COLORS"))) }}{{ $v }}{{ end }}' | COLORS=red,green,red temple
    redgreen
    $ echo '{{ range $v := env "COLORS" | split "," | uniq }}{{ $v }}{{ end }}' | COLORS=red,green,red temple
    redgreen
    ```

- `replace FROM TO STRING` - proxy for [strings.ReplaceAll](https://pkg.go.dev/strings#ReplaceAll) but with source string as last argument

- `upper STRING` - proxy for [strings.ToUpper](https://pkg.go.dev/strings#ToUpper)

- `lower STRING` - proxy for [strings.ToLower](https://pkg.go.dev/strings#ToLower)

- `istrue STRING` - returns true if string is non-empty and starts with "1", "t", "y" or "on".

    ```sh
    $ TPL='{{ if (istrue (env "FEATURE")) }}Enabled{{ else }}Disabled{{ end }}'
    $ echo "$TPL" | FEATURE=y temple
    Enabled
    $ echo "$TPL" | FEATURE=what temple
    Disabled
    ```

- `error STRING` - prints error and exists with non-zero code

    ```sh
    $ echo '{{ error "Dafuq!" }}' | temple
    ERROR: Dafuq!
    ```
