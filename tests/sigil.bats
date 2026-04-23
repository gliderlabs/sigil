setup() {
  GOOS=$(go env GOOS)
  GOARCH=$(go env GOARCH)
  if [[ "$GOOS" == "darwin" ]]; then
    export SIGIL="${SIGIL:-build/${GOOS}/gliderlabs-sigil}"
  else
    export SIGIL="${SIGIL:-build/${GOOS}/gliderlabs-sigil-${GOARCH}}"
  fi
}

@test "POSIX variable substitution" {
  result=$(echo 'Hello, $name' | $SIGIL -p name=Jeff)
  [[ "$result" == "Hello, Jeff" ]]
}

@test "POSIX variable with default value" {
  result=$(echo 'Hello, ${name:-Jeff}' | $SIGIL -p)
  [[ "$result" == "Hello, Jeff" ]]
}

@test "POSIX variable check with value set" {
  result=$(echo 'Hello, ${name:?}' | $SIGIL -p name=Jeff)
  [[ "$result" == "Hello, Jeff" ]]
}

@test "POSIX variable check fails when unset" {
  run bash -c "echo 'Hello, \${name:?}' | $SIGIL -p"
  [[ "$status" -ne 0 ]]
}

@test "template variable substitution" {
  result=$(echo 'Hello, {{ $name }}' | $SIGIL name=Jeff)
  [[ "$result" == "Hello, Jeff" ]]
}

@test "range with stdin" {
  result=$(echo '${name} is{{ range seq ${count:-3} }} cool{{ end }}!' | $SIGIL -p name=Sigil)
  [[ "$result" == "Sigil is cool cool cool!" ]]
}

@test "range variable with stdin" {
  result=$(echo 'Sigil is{{ range $i := seq 3 }} cool{{ end }}!' | $SIGIL)
  [[ "$result" == "Sigil is cool cool cool!" ]]
}

@test "range variable inline with equals" {
  result=$($SIGIL -i='Sigil is{{ range $i := seq 3 }} cool{{ end }}!')
  [[ "$result" == "Sigil is cool cool cool!" ]]
}

@test "range variable inline without equals" {
  result=$($SIGIL -i 'Sigil is{{ range $i := seq 3 }} cool{{ end }}!')
  [[ "$result" == "Sigil is cool cool cool!" ]]
}

@test "capitalize function" {
  result=$(echo 'hello {{capitalize "jeff"}}' | $SIGIL)
  [[ "$result" == "hello Jeff" ]]
}

@test "exists returns true for relative path" {
  result=$(echo '{{exists "Makefile"}}' | $SIGIL)
  [[ "$result" == "true" ]]
}

@test "exists returns true for full path" {
  result=$(echo "{{exists \"$(pwd)/Makefile\"}}" | $SIGIL)
  [[ "$result" == "true" ]]
}

@test "exists returns false for missing relative path" {
  result=$(echo '{{exists "FileNotExist"}}' | $SIGIL)
  [[ "$result" == "false" ]]
}

@test "exists returns false for missing full path" {
  result=$(echo "{{exists \"$(pwd)/FileNotExist\"}}" | $SIGIL)
  [[ "$result" == "false" ]]
}

@test "literal text passthrough" {
  result=$(echo 'XXX' | $SIGIL)
  [[ "$result" == "XXX" ]]
}

@test "split and join pipeline" {
  result=$(echo 'one,two,three' | $SIGIL -i '{{ stdin | split "," | join ":" }}')
  [[ "$result" == "one:two:three" ]]
}

@test "splitkv and joinkv pipeline" {
  result=$(echo -n 'one:two,three:four' | $SIGIL -i '{{ stdin | split "," | splitkv ":" | joinkv "=" | join "," }}')
  [[ "$result" == "one=two,three=four" || "$result" == "three=four,one=two" ]]
}

@test "JSON parse and serialize" {
  result=$(echo '{"one": "two"}' | $SIGIL -i '{{ stdin | json | tojson }}')
  [[ "$result" == "{\"one\":\"two\"}" ]]
}

@test "JSON deep nested parse and serialize" {
  result=$(echo '{"foo": {"one": "two"}}' | $SIGIL -i '{{ stdin | json | tojson }}')
  [[ "$result" == '{"foo":{"one":"two"}}' ]]
}

@test "YAML parse and serialize" {
  yaml="$(echo -e "one: two\nthree:\n- four\n- five")"
  result="$(echo -e "$yaml" | $SIGIL -i '{{ stdin | yaml | toyaml }}')"
  [[ "$result" == "$yaml" ]]
}

@test "shell command execution" {
  result="$($SIGIL -i '{{ sh "date +%m-%d-%Y" }}')"
  [[ "$result" == "$(date +%m-%d-%Y)" ]]
}

@test "HTTP GET request" {
  result="$($SIGIL -i '{{ httpget "https://httpbin.org/get" | json | pointer "/url" }}')"
  [[ "$result" == "https://httpbin.org/get" ]]
}

@test "custom delimiters" {
  result="$(SIGIL_DELIMS={{{,}}} $SIGIL -i '{{ hello {{{ $name }}} }}' name=packer)"
  [[ "$result" == "{{ hello packer }}" ]]
}

@test "substring with range" {
  result="$($SIGIL -i '{{ "abcdefgh" | substr "1:4" }}')"
  [[ "$result" == "bcd" ]]
}

@test "substring with single index" {
  result="$($SIGIL -i '{{ "abcdefgh" | substr ":4" }}')"
  [[ "$result" == "abcd" ]]
}

@test "YAML to JSON conversion" {
  result="$(printf 'joe:\n  age: 32\n  color: red' | $SIGIL -i '{{ stdin |  yaml | tojson }}')"
  [[ "$result" == '{"joe":{"age":32,"color":"red"}}' ]]
}

@test "YAML to JSON deep conversion" {
  result="$(
    $SIGIL -i '{{ stdin |  yaml | tojson }}' <<EOF
a: Easy!
b:
  c: 2
  d:
  - 3
  - 4
c:
  list:
  - one
  - two
  - tree
EOF
  )"
  [[ "$result" == '{"a":"Easy!","b":{"c":2,"d":[3,4]},"c":{"list":["one","two","tree"]}}' ]]
}

@test "JMESPath query" {
  result="$(echo '[{"name":"bob","age":20},{"name":"jim","age":30},{"name":"joe","age":40}]' | $SIGIL -i '{{stdin | json | jmespath "[? age >= `30`].name | reverse(@)"  | join ","}}')"
  [[ "$result" == 'joe,jim' ]]
}

@test "base64 encode" {
  result="$(echo 'happybirthday' | $SIGIL -i '{{ stdin | base64enc }}')"
  [[ "$result" == "aGFwcHliaXJ0aGRheQo=" ]]
}

@test "base64 decode" {
  result="$(echo 'aGFwcHliaXJ0aGRheQo=' | $SIGIL -i '{{ stdin | base64dec }}')"
  [[ "$result" == "happybirthday" ]]
}

@test "in-place basic substitution" {
  tmpfile=$(mktemp)
  echo 'Hello, {{ $name }}' >"$tmpfile"
  $SIGIL --in-place -f "$tmpfile" name=Jeff
  result=$(cat "$tmpfile")
  rm -f "$tmpfile"
  [[ "$result" == "Hello, Jeff" ]]
}

@test "in-place preserves file permissions" {
  tmpfile=$(mktemp)
  echo 'Hello, {{ $name }}' >"$tmpfile"
  chmod 0755 "$tmpfile"
  $SIGIL --in-place -f "$tmpfile" name=Jeff
  perms=$(stat -c %a "$tmpfile" 2>/dev/null || stat -f %Lp "$tmpfile")
  rm -f "$tmpfile"
  [[ "$perms" == "755" ]]
}

@test "in-place POSIX mode" {
  tmpfile=$(mktemp)
  echo 'Hello, $name' >"$tmpfile"
  $SIGIL --in-place -p -f "$tmpfile" name=Jeff
  result=$(cat "$tmpfile")
  rm -f "$tmpfile"
  [[ "$result" == "Hello, Jeff" ]]
}

@test "in-place requires filename" {
  run $SIGIL --in-place -i 'Hello'
  [[ "$status" -ne 0 ]]
}

@test "in-place rejects stdin" {
  run bash -c "echo 'Hello' | $SIGIL --in-place"
  [[ "$status" -ne 0 ]]
}

@test "in-place error preserves original file" {
  tmpfile=$(mktemp)
  echo 'Hello, {{ $name }' >"$tmpfile"
  original=$(cat "$tmpfile")
  run $SIGIL --in-place -f "$tmpfile" name=Jeff
  result=$(cat "$tmpfile")
  rm -f "$tmpfile"
  [[ "$result" == "$original" ]]
}

@test "in-place with long filename flag" {
  tmpfile=$(mktemp)
  echo 'Hello, {{ $name }}' >"$tmpfile"
  $SIGIL --in-place --filename "$tmpfile" name=World
  result=$(cat "$tmpfile")
  rm -f "$tmpfile"
  [[ "$result" == "Hello, World" ]]
}
