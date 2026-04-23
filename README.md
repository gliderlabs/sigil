# Sigil

Standalone string interpolator and template processor

```shell
echo '${name} is{{ range seq ${count:-3} }} cool{{ end }}!' | sigil -p name=Sigil
```

```text
Sigil is cool cool cool!
```

Sigil is a command line tool for template processing and POSIX-compliant
variable expansion. It was created for configuration templating, but can be used
for any text processing.

## Getting Sigil

```shell
curl -L "https://github.com/gliderlabs/sigil/releases/download/v0.10.0/gliderlabs-sigil_0.10.0_$(uname -sm|tr \  _).tgz" \
    | tar -zxC /usr/local/bin
```

Other releases can be downloaded from [Github Releases](https://github.com/gliderlabs/sigil/releases).

## Using Sigil

Template text can be provided via STDIN or from a file if provided with the `-f`
flag. Any other arguments are key-values in the form `<key>=<value>`. They are
used as variables.

* `echo 'Hello, $name' | sigil -p name=Jeff`
* `sigil -p -f config.tmpl var1=foo "var2=Hello world"`

### Variables

#### POSIX style

There are two forms of variable syntax in Sigil. The first is POSIX style, which
among other features allows default values or enforces required values:

* `$variable` - normal POSIX style
* `${variable:-"default"}` - expansion with default value
* `${variable:?}` - fails when not set

Environment variables are also available as POSIX style variables. This makes
Sigil great for quick and simple string interpolation.

#### Template style

The other syntax to use variables is consistent with the rest of the templating
syntax. It uses `{{` and `}}` to define template expressions. Variable expansion
in this form is simply used as:

* `{{ $variable }}`

You can do much more with this syntax, such as modifier pipelines. All of which
is explained below.

#### Custom Delimiters

Sometimes you want to use sigil to generate text, which uses golang templating itself.
For example if you want to generate [packer](https://www.packer.io/docs/) configuration
your template might contain a lot of `{{` and `}}`.

Instead of replacing all `{{` with `{{“{{”}}`, you can change the delimiters,
by setting the `SIGIL_DELIMS` environment variable. It is the left and right
delimiter strings, separated by a coma.

```shell
SIGIL_DELIMS={{{,}}}  sigil -i 'hello {{{ $name }}}' name=packer
```

### In-place editing

Use the `--in-place` flag with `-f` to write the rendered output back to the
template file, similar to `sed -i`:

```shell
sigil --in-place -f config.tmpl var1=foo var2=bar
```

This safely replaces the file using an atomic write (temp file + rename), so the
original file is preserved if template processing fails. File permissions are
retained.

Note: `--in-place` requires the `-f` flag. It cannot be used with `-i` (inline)
or stdin input.

### Variables from files

Use the `--vars-file` (or `-V`) flag to load variables from a file instead of
passing them all as command-line arguments. The flag can be specified multiple
times; files are merged in order, and CLI `key=value` arguments override any
file-sourced variables.

The file format is auto-detected by extension:

* `.json` - parsed as a JSON object
* `.yaml` / `.yml` - parsed as a YAML mapping
* `.env` or any other extension - parsed as key=value lines

#### JSON example (`vars.json`)

```json
{
  "name": "Jeff",
  "greeting": "Hello",
  "port": 8080
}
```

#### YAML example (`vars.yaml`)

```yaml
name: Jeff
greeting: Hello
port: 8080
```

#### Env example (`vars.env`)

```shell
# Database config
DB_HOST=localhost
DB_PORT=5432
DB_NAME="my_database"
export APP_ENV='production'
```

#### Usage

```shell
# JSON vars file
sigil -f config.tmpl -V vars.json

# YAML vars file
sigil -f config.tmpl -V vars.yaml

# env-style vars file
sigil -f config.tmpl -V vars.env

# Multiple files, later overrides earlier
sigil -f config.tmpl -V defaults.yaml -V overrides.json

# CLI args override file vars
sigil -f config.tmpl -V vars.json name=override
```

### Functions

There are a number of builtin functions that can be used as modifiers,
conditional tests, expansion data sources, and more. There are two references
for functions available:

* [Sigil builtins](http://godoc.org/github.com/gliderlabs/sigil/builtin)
* [Go template builtins](http://golang.org/pkg/text/template/#hdr-Functions)

Here are a few examples:

* `{{ $variable | capitalize }}`
* `{{ $variable | default "fallback" }}` - use fallback when variable is not provided or empty
* `{{ include "file.tmpl" "var1=foo" "var2=bar" }}`
* `{{ file "example.txt" | replace "old" "new" }}`
* `{{ json "file.json" | pointer "/Widgets/0/Name" }}`

### Conditionals

* `{{ if expr }} true {{ end }}`
* `{{ if expr }} true {{ else }} false {{ end }}`
* `{{ if expr }} true {{ else if expr }} also true {{ end }}`

### Loops / Iteration

* `{{ range expr }} element: {{.}} {{ end }}`
* `{{ range expr }} elements {{ else }} no elements {{ end }}`

### Full Syntax

Lots more is possible with this template syntax. Sigil is based on Go's
[text/template package](http://golang.org/pkg/text/template/). You can read full
documentation there.

## License

BSD
![beacon](https://ga-beacon.appspot.com/UA-58928488-2/sigil/readme?pixel "beacon")
