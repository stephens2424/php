<?php

$l = 1;
$s = "";
$a = array();

tidy_repair_string($s, $l, $l, $l);
tidy_repair_string($s, $s, $s, $s);
tidy_repair_string($l, $l, $l ,$l);
tidy_repair_string($a, $a, $a, $a);

tidy_repair_file($s, $l, $l, $l);
tidy_repair_file($s, $s, $s, $s);
tidy_repair_file($l, $l, $l ,$l);
tidy_repair_file($a, $a, $a, $a);

echo "Done\n";
?>
