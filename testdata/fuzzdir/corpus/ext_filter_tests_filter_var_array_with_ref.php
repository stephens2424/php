<?php

$array = ["123foo"];
$array2 = [&$array];
var_dump(filter_var_array($array2, FILTER_VALIDATE_INT));
var_dump($array);

?>
