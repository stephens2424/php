<?php

$array = ['a', 'b'];
$ref =& $array[0];

var_dump(array_count_values($array));

?>
