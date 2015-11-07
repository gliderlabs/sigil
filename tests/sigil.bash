GOOS=$(go env GOOS)
export SIGIL="${SIGIL:-build/${GOOS^}/sigil}"

T_entries_with_join() {
    result=$(echo '{{ $y := yaml "tests/test.yml"}}{{ $y.images | entries "%v=%v" | join ","  }}' | $SIGIL)
    echo "HACK=$result"
    [[ "$result" == "us-east-1=ami-11111111,us-west-1=ami-22222222" ]]
}

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
