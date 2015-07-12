<?php

$php = getenv('TEST_PHP_EXECUTABLE');

var_dump(`"$php" -nd max_execution_time=111 -r 'var_dump(ini_get("max_execution_time"));'`);
var_dump(`"$php" -nd max_execution_time=500 -r 'var_dump(ini_get("max_execution_time"));'`);

?>
===DONE===
