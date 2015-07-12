<?php

	require_once("connect.inc");

	/*** test mysqli_connect 127.0.0.1 ***/
	$link = my_mysqli_connect($host, $user, $passwd, $db, $port, $socket);
	mysqli_select_db($link, $db);

	mysqli_query($link, "DROP TABLE IF EXISTS test_warnings");
	mysqli_query($link, "DROP TABLE IF EXISTS test_warnings");

	var_dump(mysqli_warning_count($link));

	mysqli_close($link);
	print "done!";
?>
