<?php

$string = '';

// These two in any order
$string .= "\r\n";
$string .= "''''";

// Total string length > 8192
$string .= str_repeat(chr(rand(32, 127)), 8184);

// Ending in this string
$string .= "say";

$finfo = new finfo();
$type = $finfo->buffer($string);
var_dump($type);

?>
