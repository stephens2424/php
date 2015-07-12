<?php

$handler = "flatfile";
require_once(dirname(__FILE__) .'/test.inc');
echo "database handler: $handler\n";

$db_file1 = dba_popen($db_filename, 'n-t', 'flatfile');

?>
===DONE===
