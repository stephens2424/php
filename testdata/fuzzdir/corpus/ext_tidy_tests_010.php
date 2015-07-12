<?php 
$a = tidy_parse_string("<HTML><BODY BGCOLOR=#FFFFFF ALINK=#000000></BODY></HTML>", array('newline' => 'LF'));
var_dump($a->root());
var_dump($a->body());
var_dump($a->html());
var_dump($a->head());

?>
