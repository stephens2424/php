<?php
$handler = "inifile";
include "test.inc";

$dba = dba_open($db_filename, "n", $handler)
	or die;
for ($i = 0; $i < 3; ++$i) {
	echo "insert $i:";
	var_dump(dba_insert("a", $i, $dba));
}

echo "exists:";
var_dump(dba_exists("a", $dba));
echo "delete:";
var_dump(dba_delete("a", $dba));
echo "exists:";
var_dump(dba_exists("a", $dba));
echo "delete:";
var_dump(dba_delete("a", $dba));

?>
===DONE===
