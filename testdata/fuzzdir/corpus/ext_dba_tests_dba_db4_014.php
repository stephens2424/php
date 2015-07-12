<?php

$handler = "db4";
require_once(dirname(__FILE__) .'/test.inc');
echo "database handler: $handler\n";

if (($db_file = dba_open($db_filename, "wl", $handler)) !== FALSE) {
    echo "database file opened\n";
    dba_close($db_file);
} else {
    echo "Error creating $db_filename\n";
}

?>
