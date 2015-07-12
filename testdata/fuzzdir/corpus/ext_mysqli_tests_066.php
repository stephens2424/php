<?php
	require_once("connect.inc");

	/*** test mysqli_connect 127.0.0.1 ***/
	$mysql = new my_mysqli($host, $user, $passwd, $db, $port, $socket);

	if (!mysqli_query($mysql, "SET sql_mode=''"))
		printf("[002] Cannot set SQL-Mode, [%d] %s\n", mysqli_errno($mysql), mysqli_error($mysql));

	$mysql->query("DROP TABLE IF EXISTS test_warnings");

	$mysql->query("CREATE TABLE test_warnings (a int not null) ENGINE=myisam");

	$mysql->query("INSERT INTO test_warnings VALUES (1),(2),(NULL)");

	if (($warning = $mysql->get_warnings())) {
		do {
			printf("Warning\n");
		} while ($warning->next());
	}

	$mysql->close();
	print "done!";
?>
