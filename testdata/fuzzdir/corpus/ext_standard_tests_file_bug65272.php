<?php

$file = dirname(__FILE__)."/flock.dat";

$fp1 = fopen($file, "w");
var_dump(flock($fp1, LOCK_SH));

$fp2 = fopen($file, "r");
var_dump(flock($fp2, LOCK_EX|LOCK_NB, $wouldblock));
var_dump($wouldblock);

@unlink($file);
echo "Done\n";
?>
