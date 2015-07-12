<?php
	require_once("connect.inc");

	/*** test mysqli_connect 127.0.0.1 ***/
	$mysql = my_mysqli_connect($host, $user, $passwd, $db, $port, $socket);

	$mysql->query("DROP TABLE IF EXISTS litest");
	$mysql->query("CREATE TABLE litest (a VARCHAR(20))");

	$mysql->query("LOAD DATA LOCAL INFILE 'filenotfound' INTO TABLE litest");
	var_dump($mysql->error);

	$mysql->close();
	printf("Done");
?>
