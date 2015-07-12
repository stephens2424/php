<?php

var_dump(filter_has_var(INPUT_GET, "a"));
var_dump(filter_has_var(INPUT_GET, "abc"));
var_dump(filter_has_var(INPUT_GET, "nonex"));
var_dump(filter_has_var(INPUT_GET, " "));
var_dump(filter_has_var(INPUT_GET, ""));
var_dump(filter_has_var(INPUT_GET, array()));

var_dump(filter_has_var(INPUT_POST, "b"));
var_dump(filter_has_var(INPUT_POST, "bbc"));
var_dump(filter_has_var(INPUT_POST, "nonex"));
var_dump(filter_has_var(INPUT_POST, " "));
var_dump(filter_has_var(INPUT_POST, ""));
var_dump(filter_has_var(INPUT_POST, array()));

var_dump(filter_has_var(-1, ""));
var_dump(filter_has_var("", ""));
var_dump(filter_has_var(array(), array()));
var_dump(filter_has_var(array(), ""));
var_dump(filter_has_var("", array()));

echo "Done\n";
?>
