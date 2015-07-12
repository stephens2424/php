<?php
$a="aaa\n<>";

var_dump( mb_ereg("^[^><]+$",$a) );
var_dump( !!preg_match("/^[^><]+$/",$a) );
?>

