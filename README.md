# `temple`

So I needed template executor for nginx configs that would go into production image.

I ended up writing this tool. It has zero dependencies. Output binary is around 2 MB.

Why not ...

- [`envsubst`](https://man7.org/linux/man-pages/man1/envsubst.1.html) - does simple substitution and does not support default values, conditions, or loops
- [`gomplate`](https://github.com/hairyhenderson/gomplate) - because it's too fat! Seriously, guys, 100 megs for a template renderer?!
- `sed`/`awk` - too much fuss

## Installation

```sh
go install github.com/vegas503/temple/cmd/temple@latest
```

## Usage

All the heavy lifting is done by Go's [text/template](https://pkg.go.dev/text/template) package. Go see their docs first!

In addition to the built-ins, `temple` provides some additional functions for common scenarios:

- `env STRING` - gets env var value, throws error if env var is not set

    ```sh
    $ echo '{{ env "USER" }}' | temple
    username
    ```

- `envdefault STRING STRING` - gets env var value, returns default value if empty or unset

    ```sh
    $ echo '{{ envdefault "SOME_UNSET_VAR" "foo" }}' | temple
    foo
    ```

- `split STRING DELIM` - splits a string with a delimiter, trims whitespaces and returns array with all empty elements removed

    ```sh
    $ echo '{{ range (split "A, B , C" ",") }}{{ . }}{{ end }}' | temple
    ABC
    ```

- `coalesce STRING...` - returns first non-empty string

    ```sh
    $ echo '{{ coalesce "" "two" "three" }}' | temple
    two
    ```

- `append ARRAY...` - concatenates arrays

    ```sh
    $ cat temple.tpl
    {{- $colors := (split (env "COLORS") ",") }}
    {{- $numbers := (split (env "NUMBERS") ",") }}
    {{- range (append $colors $numbers) }}
      * {{ . }}
    {{- end }}
    $ COLORS=red,green,blue NUMBERS=34,42 temple -i temple.tpl
      * red
      * green
      * blue
      * 34
      * 42

- `uniq ARRAY` - returns array with all duplicate elements removed

    ```sh
    $ echo '{{ range (uniq (split (env "COLORS") ",")) }}{{ . }}{{ end }}' | COLORS=red,green,red temple
    redgreen

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
