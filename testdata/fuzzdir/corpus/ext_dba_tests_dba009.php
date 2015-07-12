<?php
	require_once(dirname(__FILE__) .'/test.inc');
	echo "database handler: $handler\n";
	if (($db=dba_popen($db_file, "n", $handler))!==FALSE) {
		echo "Opened\n";
		dba_insert("a", "Inserted", $db);
		echo dba_fetch("a", $db)."\n";
		dba_close($db);
		echo "Closed\n";
	} else {
		echo "Error creating database\n";
	}
	if (($db=dba_popen($db_file, "n", $handler))!==FALSE) {
		echo "Opened\n";
		dba_insert("a", "Inserted", $db);
		echo dba_fetch("a", $db)."\n";
	}
?>
