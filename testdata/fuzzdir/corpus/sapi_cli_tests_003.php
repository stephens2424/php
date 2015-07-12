<?php

$php = getenv('TEST_PHP_EXECUTABLE');

var_dump(`$php -n -d max_execution_time=111 -r 'var_dump(ini_get("max_execution_time"));'`);
var_dump(`$php -n -d max_execution_time=500 -r 'var_dump(ini_get("max_execution_time"));'`);
var_dump(`$php -n -d max_execution_time=500 -d max_execution_time=555 -r 'var_dump(ini_get("max_execution_time"));'`);
var_dump(`$php -n -d upload_tmp_dir=/test/path -d max_execution_time=555 -r 'var_dump(ini_get("max_execution_time")); var_dump(ini_get("upload_tmp_dir"));'`);

echo "Done\n";
?>
