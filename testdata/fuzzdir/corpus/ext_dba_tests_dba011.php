<?php
require_once(dirname(__FILE__) .'/test.inc');
echo "database handler: $handler\n";
var_dump(dba_open($db_file));
var_dump(dba_open($db_file, 'n'));
var_dump(dba_open($db_file, 'n', 'bogus'));
var_dump(dba_open($db_file, 'q', $handler));
var_dump(dba_open($db_file, 'nq', $handler));
var_dump(dba_open($db_file, 'n', $handler, 2, 3, 4, 5, 6, 7, 8));
?>
