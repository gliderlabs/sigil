# shellcheck disable=all
GOOS=$(go env GOOS)
export SIGIL="${SIGIL:-build/${GOOS}/gliderlabs-sigil}-amd64"

T_posix_var() {
  result=$(echo 'Hello, $name' | $SIGIL -p name=Jeff)
  [[ "$result" == "Hello, Jeff" ]]
}

T_posix_var_default() {
  result=$(echo 'Hello, ${name:-Jeff}' | $SIGIL -p)
  [[ "$result" == "Hello, Jeff" ]]
}

T_posix_var_check() {
  result=$(echo 'Hello, ${name:?}' | $SIGIL -p name=Jeff)
  [[ "$result" == "Hello, Jeff" ]]
}

T_posix_var_check_unset() {
  echo 'Hello, ${name:?}' | $SIGIL -p &>/dev/null
  [[ $? -ne 0 ]]
}

T_template_var() {
  result=$(echo 'Hello, {{ $name }}' | $SIGIL name=Jeff)
  [[ "$result" == "Hello, Jeff" ]]
}

T_range_stdin() {
  result=$(echo '${name} is{{ range seq ${count:-3} }} cool{{ end }}!' | $SIGIL -p name=Sigil)
  [[ "$result" == "Sigil is cool cool cool!" ]]
}

T_range_var_stdin() {
  result=$(echo 'Sigil is{{ range $i := seq 3 }} cool{{ end }}!' | $SIGIL)
  [[ "$result" == "Sigil is cool cool cool!" ]]
}

T_range_var_inline_with_equals() {
  result=$($SIGIL -i='Sigil is{{ range $i := seq 3 }} cool{{ end }}!')
  [[ "$result" == "Sigil is cool cool cool!" ]]
}

T_range_var_inline_without_equals() {
  result=$($SIGIL -i 'Sigil is{{ range $i := seq 3 }} cool{{ end }}!')
  [[ "$result" == "Sigil is cool cool cool!" ]]
}

T_capitalize() {
  result=$(echo 'hello {{capitalize "jeff"}}' | $SIGIL)
  [[ "$result" == "hello Jeff" ]]
}

T_exists_exist_relative() {
  result=$(echo '{{exists "Makefile"}}' | $SIGIL)
  [[ "$result" == "true" ]]
}

T_exists_exist_full() {
  result=$(echo "{{exists \"$(pwd)/Makefile\"}}" | $SIGIL)
  [[ "$result" == "true" ]]
}

T_exists_notexist_relative() {
  result=$(echo '{{exists "FileNotExist"}}' | $SIGIL)
  [[ "$result" == "false" ]]
}

T_exists_notexist_full() {
  result=$(echo "{{exists \"$(pwd)/FileNotExist\"}}" | $SIGIL)
  [[ "$result" == "false" ]]
}

T_XXX() {
  result=$(echo 'XXX' | $SIGIL)
  [[ "$result" == "XXX" ]]
}

T_split_join() {
  result=$(echo 'one,two,three' | $SIGIL -i '{{ stdin | split "," | join ":" }}')
  [[ "$result" == "one:two:three" ]]
}

T_splitkv_joinkv() {
  result=$(echo -n 'one:two,three:four' | $SIGIL -i '{{ stdin | split "," | splitkv ":" | joinkv "=" | join "," }}')
  [[ "$result" == "one=two,three=four" || "$result" == "three=four,one=two" ]]
}

T_json() {
  result=$(echo '{"one": "two"}' | $SIGIL -i '{{ stdin | json | tojson }}')
  [[ "$result" == "{\"one\":\"two\"}" ]]
}

T_json_deep() {
  result=$(echo '{"foo": {"one": "two"}}' | $SIGIL -i '{{ stdin | json | tojson }}')
  [[ "$result" == '{"foo":{"one":"two"}}' ]]
}

T_yaml() {
  yaml="$(echo -e "one: two\nthree:\n- four\n- five")"
  result="$(echo -e "$yaml" | $SIGIL -i '{{ stdin | yaml | toyaml }}')"
  [[ "$result" == "$yaml" ]]
}

T_shell() {
  result="$($SIGIL -i '{{ sh "date +%m-%d-%Y" }}')"
  [[ "$result" == "$(date +%m-%d-%Y)" ]]
}

T_httpget() {
  result="$($SIGIL -i '{{ httpget "https://httpbin.org/get" | json | pointer "/url" }}')"
  [[ "$result" == "https://httpbin.org/get" ]]
}

T_custom_delim() {
  result="$(SIGIL_DELIMS={{{,}}} $SIGIL -i '{{ hello {{{ $name }}} }}' name=packer)"
  [[ "$result" == "{{ hello packer }}" ]]
}

T_substr() {
  result="$($SIGIL -i '{{ "abcdefgh" | substr "1:4" }}')"
  [[ "$result" == "bcd" ]]
}
T_substr_single_index() {
  result="$($SIGIL -i '{{ "abcdefgh" | substr ":4" }}')"
  [[ "$result" == "abcd" ]]
}

T_yamltojson() {
  result="$(printf 'joe:\n  age: 32\n  color: red' | $SIGIL -i '{{ stdin |  yaml | tojson }}')"
  [[ "$result" == '{"joe":{"age":32,"color":"red"}}' ]]
}

T_yamltojsondeep() {
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

T_jmespath() {
  result="$(echo '[{"name":"bob","age":20},{"name":"jim","age":30},{"name":"joe","age":40}]' | $SIGIL -i '{{stdin | json | jmespath "[? age >= `30`].name | reverse(@)"  | join ","}}')"
  [[ "$result" == 'joe,jim' ]]
}

T_base64enc() {
  result="$(echo 'happybirthday' | $SIGIL -i '{{ stdin | base64enc }}')"
  [[ "$result" == "aGFwcHliaXJ0aGRheQo=" ]]
}

T_base64dec() {
  result="$(echo 'aGFwcHliaXJ0aGRheQo=' | $SIGIL -i '{{ stdin | base64dec }}')"
  [[ "$result" == "happybirthday" ]]
}

T_inplace_basic() {
  tmpfile=$(mktemp)
  echo 'Hello, {{ $name }}' >"$tmpfile"
  $SIGIL --in-place -f "$tmpfile" name=Jeff
  result=$(cat "$tmpfile")
  rm -f "$tmpfile"
  [[ "$result" == "Hello, Jeff" ]]
}

T_inplace_permissions() {
  tmpfile=$(mktemp)
  echo 'Hello, {{ $name }}' >"$tmpfile"
  chmod 0755 "$tmpfile"
  $SIGIL --in-place -f "$tmpfile" name=Jeff
  perms=$(stat -c %a "$tmpfile" 2>/dev/null || stat -f %Lp "$tmpfile")
  rm -f "$tmpfile"
  [[ "$perms" == "755" ]]
}

T_inplace_posix() {
  tmpfile=$(mktemp)
  echo 'Hello, $name' >"$tmpfile"
  $SIGIL --in-place -p -f "$tmpfile" name=Jeff
  result=$(cat "$tmpfile")
  rm -f "$tmpfile"
  [[ "$result" == "Hello, Jeff" ]]
}

T_inplace_requires_filename() {
  $SIGIL --in-place -i 'Hello' 2>/dev/null
  [[ $? -ne 0 ]]
}

T_inplace_rejects_stdin() {
  echo 'Hello' | $SIGIL --in-place 2>/dev/null
  [[ $? -ne 0 ]]
}

T_inplace_error_preserves_file() {
  tmpfile=$(mktemp)
  echo 'Hello, {{ $name }' >"$tmpfile"
  original=$(cat "$tmpfile")
  $SIGIL --in-place -f "$tmpfile" name=Jeff 2>/dev/null
  result=$(cat "$tmpfile")
  rm -f "$tmpfile"
  [[ "$result" == "$original" ]]
}

T_inplace_long_flag() {
  tmpfile=$(mktemp)
  echo 'Hello, {{ $name }}' >"$tmpfile"
  $SIGIL --in-place --filename "$tmpfile" name=World
  result=$(cat "$tmpfile")
  rm -f "$tmpfile"
  [[ "$result" == "Hello, World" ]]
}
