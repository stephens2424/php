<?php

$n = "10";
$n .= "0";
$nums = [&$n, 100];
var_dump(array_sum($nums));
var_dump($n);

?>
