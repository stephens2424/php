<?php
// Sinhala characters
var_dump(preg_replace('/(?<!ක)/u', '*', 'ක'));
var_dump(preg_replace('/(?<!ක)/u', '*', 'ම'));
// English characters
var_dump(preg_replace('/(?<!k)/u', '*', 'k'));
var_dump(preg_replace('/(?<!k)/u', '*', 'm'));
// Sinhala characters
preg_match_all('/(?<!ක)/u', 'ම', $matches, PREG_OFFSET_CAPTURE);
var_dump($matches);
// invalid UTF-8
var_dump(preg_replace('/(?<!ක)/u', '*', "\xFCක"));
var_dump(preg_replace('/(?<!ක)/u', '*', "ක\xFC"));
var_dump(preg_match_all('/(?<!ක)/u', "\xFCම", $matches, PREG_OFFSET_CAPTURE));
var_dump(preg_match_all('/(?<!ක)/u', "\xFCම", $matches, PREG_OFFSET_CAPTURE));
?>
