<?php

$key = ftok(__FILE__, 't');
$s = shm_attach($key, 1024);

shm_put_var($s, -1, "test string");
shm_put_var($s, 0, new stdclass);
shm_put_var($s, 1, array(1,2,3));
shm_put_var($s, 2, false);
shm_put_var($s, 3, null);

var_dump(shm_get_var());

var_dump(shm_get_var(-1, -1));

var_dump(shm_get_var($s, 1000));
var_dump(shm_get_var($s, -10000));

var_dump(shm_get_var($s, array()));
var_dump(shm_get_var($s, -1));
var_dump(shm_get_var($s, 0));
var_dump(shm_get_var($s, 1));
var_dump(shm_get_var($s, 2));
var_dump(shm_get_var($s, 3));

shm_put_var($s, 3, "test");
shm_put_var($s, 3, 1);
shm_put_var($s, 3, null);

var_dump(shm_get_var($s, 3));
shm_remove($s);

echo "Done\n";
?>
