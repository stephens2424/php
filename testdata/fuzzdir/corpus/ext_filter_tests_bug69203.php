<?php
var_dump(filter_var("\x7f", FILTER_SANITIZE_STRING, FILTER_FLAG_STRIP_HIGH));
var_dump(filter_var("\x7f", FILTER_UNSAFE_RAW, FILTER_FLAG_STRIP_HIGH));
var_dump(filter_var("\x7f", FILTER_SANITIZE_ENCODED, FILTER_FLAG_STRIP_HIGH));
var_dump(filter_var("\x7f", FILTER_SANITIZE_SPECIAL_CHARS, FILTER_FLAG_STRIP_HIGH));
?>
