<?php

error_reporting(E_ALL);

require_once(dirname(__FILE__) .'/test.inc');

$db = dba_popen($db_filename, 'c');

$key = 1;
$copy = $key;

echo gettype($key)."\n";
echo gettype($copy)."\n";

dba_exists($key, $db);

echo gettype($key)."\n";
echo gettype($copy)."\n";

dba_close($db);

?>
