<?php

var_dump(json_decode('"\ud834"'));
var_dump(json_last_error() === JSON_ERROR_UTF16);
var_dump(json_last_error_msg());
?>
