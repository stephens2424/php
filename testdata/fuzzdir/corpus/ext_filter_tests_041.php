<?php 
var_dump($_COOKIE);
var_dump(filter_has_var(INPUT_COOKIE, "abc"));
var_dump(filter_input(INPUT_COOKIE, "abc"));
var_dump(filter_input(INPUT_COOKIE, "def"));
var_dump(filter_input(INPUT_COOKIE, "xyz"));
var_dump(filter_has_var(INPUT_COOKIE, "bogus"));
var_dump(filter_input(INPUT_COOKIE, "xyz", FILTER_SANITIZE_SPECIAL_CHARS));
?>
