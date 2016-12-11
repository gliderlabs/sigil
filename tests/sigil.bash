GOOS=$(go env GOOS)
export SIGIL="${SIGIL:-build/${GOOS^}/sigil}"

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
  echo 'Hello, ${name:?}' | $SIGIL -p &> /dev/null
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

T_exists() {
  result=$(echo '{{exists "Makefile"}}' | $SIGIL)
  [[ "$result" == "true" ]]
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

