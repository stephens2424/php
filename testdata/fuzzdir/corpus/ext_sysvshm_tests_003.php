<?php

$key = ftok(dirname(__FILE__)."/003.phpt", 'q');

var_dump(shm_detach());
var_dump(shm_detach(1,1));

$s = shm_attach($key);

var_dump(shm_detach($s));
var_dump(shm_detach($s));
shm_remove($s);

var_dump(shm_detach(0));
var_dump(shm_detach(1));
var_dump(shm_detach(-1));

echo "Done\n";
?>
