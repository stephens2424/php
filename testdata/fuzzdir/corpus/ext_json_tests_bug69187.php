<?php
var_dump(json_decode(NULL)); 
var_dump(json_last_error());
var_dump(json_decode(FALSE));
var_dump(json_last_error());
var_dump(json_decode(""));
var_dump(json_last_error());

var_dump(json_decode(0));
var_dump(json_last_error());
var_dump(json_decode(1));
var_dump(json_last_error());
var_dump(json_decode(TRUE));
var_dump(json_last_error());

json_decode("\xED\xA0\xB4");
var_dump(json_last_error());

json_decode("\x00");
var_dump(json_last_error());

json_decode("\"\xED\xA0\xB4\"");
var_dump(json_last_error());

json_decode("\"\x00\"");
var_dump(json_last_error());
?>
