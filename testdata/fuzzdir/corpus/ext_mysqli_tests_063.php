<?php
	require_once("connect.inc");

	$mysql = new my_mysqli($host, $user, $passwd, $db, $port, $socket);

	$stmt = new mysqli_stmt($mysql, "SELECT 'foo' FROM DUAL");
	$stmt->execute();
	$stmt->bind_result($foo);
	$stmt->fetch();
	$stmt->close();
	$mysql->close();

	var_dump($foo);
?>
