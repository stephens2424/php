<?php

$handler = 'db4';
require_once(dirname(__FILE__) .'/test.inc');

$db = dba_popen($db_filename, 'c', 'db4');

dba_insert('X', 'XYZ', $db);
dba_insert('Y', '123', $db);

var_dump($db, dba_fetch('X', $db));

var_dump(dba_firstkey($db));
var_dump(dba_nextkey($db));

dba_close($db);

?>
===DONE===
