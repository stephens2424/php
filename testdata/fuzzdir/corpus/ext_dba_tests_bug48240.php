<?php

$handler = 'db4';
require_once(dirname(__FILE__) .'/test.inc');

$db = dba_open($db_filename, 'c', 'db4');

var_dump(dba_nextkey($db));

dba_close($db);

?>
===DONE===
