<?php

$key = ftok(__FILE__, 't');
$s = shm_attach($key, 1024);

var_dump(shm_put_var());
var_dump(shm_put_var(-1, -1, -1));
var_dump(shm_put_var(-1, 10, "qwerty"));
var_dump(shm_put_var($s, -1, "qwerty"));
var_dump(shm_put_var($s, 10, "qwerty"));
var_dump(shm_put_var($s, 10, "qwerty"));

$string = str_repeat("test", 512);
var_dump(shm_put_var($s, 11, $string));

shm_remove($s);

echo "Done\n";
?>
