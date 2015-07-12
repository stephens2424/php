<?php

$foo = 'foo';
$arr = ['bar' => &$foo]; 
var_dump(strtr('foobar', $arr));

?>
