<?php

var_dump(preg_replace_callback_array());
var_dump(preg_replace_callback_array(1));
var_dump(preg_replace_callback_array(1,2));
var_dump(preg_replace_callback_array(1,2,3));
$a = 5;
var_dump(preg_replace_callback_array(1,2,3,$a));
$a = "";
var_dump(preg_replace_callback_array(array("" => ""),"","",$a));
$a = array();
$b = "";
var_dump(preg_replace_callback($a, $a, $a, $a, $b));
var_dump($b);
$b = "";
var_dump(preg_replace_callback_array(array("xx" => "s"), $a, -1, $b));
var_dump($b);
function f() {
	static $count = 1;
	throw new Exception($count);
}

var_dump(preg_replace_callback_array(array('/\w' => 'f'), 'z'));

try {
	var_dump(preg_replace_callback_array(array('/\w/' => 'f', '/.*/' => 'f'), 'z'));
} catch (Exception $e) {
	var_dump($e->getMessage());
}

echo "Done\n";
?>
