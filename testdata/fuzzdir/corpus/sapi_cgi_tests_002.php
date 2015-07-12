<?php
include "include.inc";

$php = get_cgi_path();
reset_env_vars();

$file = dirname(__FILE__)."/002.test.php";

file_put_contents($file, '<?php var_dump(ini_get("max_execution_time")); ?>');

var_dump(`$php -n -d max_execution_time=111 $file`);
var_dump(`$php -n -d max_execution_time=500 $file`);
var_dump(`$php -n -d max_execution_time=500 -d max_execution_time=555 $file`);

file_put_contents($file, '<?php var_dump(ini_get("max_execution_time")); var_dump(ini_get("upload_tmp_dir")); ?>');

var_dump(`$php -n -d upload_tmp_dir=/test/path -d max_execution_time=555 $file`);

unlink($file);

echo "Done\n";
?>
