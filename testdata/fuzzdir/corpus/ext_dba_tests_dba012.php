<?php
$handler = "flatfile";
require_once(dirname(__FILE__) .'/test.inc');
echo "database handler: $handler\n";

echo "Test 1\n";

ini_set('dba.default_handler', 'does_not_exist');

var_dump(dba_open($db_filename, 'c'));

echo "Test 2\n";

ini_set('dba.default_handler', '');

var_dump(dba_open($db_filename, 'n'));

?>
