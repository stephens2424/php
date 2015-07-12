<?php

$key = ftok(__FILE__, 't');
$s = shm_attach($key, 1024);

shm_put_var($s, 1, "test string");

var_dump(shm_remove_var());
var_dump(shm_remove_var(-1, -1));
var_dump(shm_remove_var($s, -10));

var_dump(shm_get_var($s, 1));

var_dump(shm_remove_var($s, 1));
var_dump(shm_get_var($s, 1));

var_dump(shm_remove_var($s, 1));
var_dump(shm_get_var($s, 1));

shm_remove($s);
echo "Done\n";
?>
