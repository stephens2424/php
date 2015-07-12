<?php
	$handler = 'tcadb';
    require_once dirname(__FILE__) .'/skipif.inc';
    $lock_flag = 'l';
    $db_filename = $db_file = dirname(__FILE__) .'/test0.tch';
    @unlink($db_filename);
    @unlink($db_filename.'.lck');
	require_once dirname(__FILE__) .'/dba_handler.inc';
?>
===DONE===
