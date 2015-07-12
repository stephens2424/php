<?php
$handler = "db4";
require_once(dirname(__FILE__) .'/test.inc');
echo "database handler: $handler\n";
$db_file1 = $db_filename1 = dirname(__FILE__).'/test1.dbm'; 
$db_file2 = $db_filename2 = dirname(__FILE__).'/test2.dbm'; 
if (($db_file=dba_open($db_file, "n", $handler))!==FALSE) {
    echo "database file created\n";
} else {
    echo "$db_file does not exist\n";
}
if (($db_file1=dba_open($db_file1, "n", $handler))!==FALSE) {
    echo "database file created\n";
} else {
    echo "$db_file does not exist\n";
}
if (($db_file2=dba_open($db_file2, "n", $handler))!==FALSE) {
    echo "database file created\n";
} else {
    echo "$db_file does not exist\n";
}
var_dump(dba_list());
dba_close($db_file);

@unlink($db_filename1);
@unlink($db_filename2);
?>
