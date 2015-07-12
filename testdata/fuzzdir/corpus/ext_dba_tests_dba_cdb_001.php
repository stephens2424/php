<?php

$handler = 'cdb';
require_once(dirname(__FILE__) .'/test.inc');

echo "Test 0\n";

if (($db_file = dba_open($db_filename, 'n', $handler))!==FALSE) {
    var_dump(dba_insert("key1", "Content String 1", $db_file));
    var_dump(dba_replace("key1", "New Content String", $db_file));
    var_dump(dba_fetch("key1", $db_file));
    var_dump(dba_firstkey($db_file));
    var_dump(dba_delete("key1", $db_file));
    var_dump(dba_optimize($db_file));
    var_dump(dba_sync($db_file));
    dba_close($db_file);
}
else {
    echo "Failed to open DB\n";
}

unlink($db_filename);

echo "Test 1\n";

if (($db_file = dba_open($db_filename, 'c', $handler))!==FALSE) {
    dba_insert("key1", "Content String 1", $db_file);
    dba_close($db_file);
}
else {
    echo "Failed to open DB\n";
}

echo "Test 2\n";

if (($db_file = dba_open($db_filename, 'r', $handler))!==FALSE) {
    dba_insert("key1", "Content String 1", $db_file);
    dba_close($db_file);
}
else {
    echo "Failed to open DB\n";
}

echo "Test 3\n";

if (($db_file = dba_open($db_filename, 'w', $handler))!==FALSE) {
    echo dba_fetch("key1", $db_file), "\n";
    dba_close($db_file);
}
else {
    echo "Failed to open DB\n";
}

?>
===DONE===
