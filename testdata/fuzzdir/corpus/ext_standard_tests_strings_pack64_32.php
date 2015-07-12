<?php
var_dump(pack("Q", 0));
var_dump(pack("J", 0));
var_dump(pack("P", 0));
var_dump(pack("q", 0));

var_dump(unpack("Q", ''));
var_dump(unpack("J", ''));
var_dump(unpack("P", ''));
var_dump(unpack("q", ''));
?>
