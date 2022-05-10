# Sigil

[![CircleCI](https://img.shields.io/circleci/project/gliderlabs/sigil/release.svg)](https://circleci.com/gh/gliderlabs/sigil)
[![IRC Channel](https://img.shields.io/badge/irc-%23gliderlabs-blue.svg)](https://kiwiirc.com/client/irc.freenode.net/#gliderlabs)


Standalone string interpolator and template processor

```
$ echo '${name} is{{ range seq ${count:-3} }} cool{{ end }}!' | sigil -p name=Sigil
Sigil is cool cool cool!
```

Sigil is a command line tool for template processing and POSIX-compliant
variable expansion. It was created for configuration templating, but can be used
for any text processing.

## Getting Sigil

```shell
curl -L "https://github.com/gliderlabs/sigil/releases/download/v0.9.0/gliderlabs-sigil_0.9.0_$(uname -sm|tr \  _).tgz" \
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

```
SIGIL_DELIMS={{{,}}}  sigil -i 'hello {{{ $name }}}' name=packer
```

### Functions

There are a number of builtin functions that can be used as modifiers,
conditional tests, expansion data sources, and more. There are two references
for functions available:

 * [Sigil builtins](http://godoc.org/github.com/gliderlabs/sigil/builtin)
 * [Go template builtins](http://golang.org/pkg/text/template/#hdr-Functions)

Here are a few examples:

 * `{{ $variable | capitalize }}`
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
<img src="https://ga-beacon.appspot.com/UA-58928488-2/sigil/readme?pixel" />
