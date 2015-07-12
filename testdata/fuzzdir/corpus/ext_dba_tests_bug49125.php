<?php

error_reporting(E_ALL);

$handler = 'db4';
require_once(dirname(__FILE__) .'/test.inc');

$db = dba_popen($db_filename, 'c', 'db4');

dba_insert('foo', 'foo', $db);

var_dump(dba_exists('foo', $db));

dba_close($db);

?>
