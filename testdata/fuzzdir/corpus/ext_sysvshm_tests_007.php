<?php

$key = ftok(__FILE__, 't');
$s = shm_attach($key, 1024);

var_dump(shm_remove());
var_dump(shm_remove(-1));
var_dump(shm_remove(0));
var_dump(shm_remove(""));

var_dump(shm_remove($s));
var_dump(shm_remove($s));

shm_detach($s);
var_dump(shm_remove($s));

echo "Done\n";
?>
