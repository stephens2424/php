<?php
	$handler = 'gdbm';
	require_once dirname(__FILE__) .'/test.inc';
	$lock_flag = ''; // lock in library
	require_once dirname(__FILE__) .'/dba_handler.inc';
	
	// Read during write is system dependent. Important is that there is no deadlock
?>
===DONE===
