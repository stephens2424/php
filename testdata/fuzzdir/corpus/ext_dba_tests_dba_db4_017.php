<?php

$handler = "db4";
require_once(dirname(__FILE__) .'/test.inc');
echo "database handler: $handler\n";

if (($db_file = dba_open($db_filename, "cd", $handler)) !== FALSE) {
    if (file_exists($db_filename)) {
        echo "database file created\n";
        var_dump(dba_insert("key1", "This is a test insert", $db_file));
        echo dba_fetch("key1", $db_file), "\n";
        dba_close($db_file);
    } else {
        echo "File did not get created\n";
    }
} else {
    echo "Error creating $db_filename\n";
}

?>
