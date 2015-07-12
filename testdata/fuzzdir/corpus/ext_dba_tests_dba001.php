<?php
	require_once(dirname(__FILE__) .'/test.inc');
	echo "database handler: $handler\n";
	if (($db_file=dba_open($db_file, "n", $handler))!==FALSE) {
    	echo "database file created\n";
		dba_close($db_file);
	} else {
    	echo "$db_file does not exist\n";
    }
?>
