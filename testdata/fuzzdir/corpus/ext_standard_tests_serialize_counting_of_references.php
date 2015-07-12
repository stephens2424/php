<?php

$ref1 = 1;
$ref2 = 2;

$arr = [&$ref1, &$ref1, &$ref2, &$ref2];
var_dump(serialize($arr));

?>
