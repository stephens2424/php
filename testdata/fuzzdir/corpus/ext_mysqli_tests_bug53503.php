<?php
	require_once("connect.inc");

	if (!$link = my_mysqli_connect($host, $user, $passwd, $db, $port, $socket)) {
		printf("[001] Connect failed, [%d] %s\n", mysqli_connect_errno(), mysqli_connect_error());
	}

	if (!$link->query("DROP TABLE IF EXISTS test")) {
		printf("[002] [%d] %s\n", $link->errno, $link->error);
	}

	if (!$link->query("CREATE TABLE test (dump1 INT UNSIGNED NOT NULL PRIMARY KEY) ENGINE=" . $engine)) {
		printf("[003] [%d] %s\n", $link->errno, $link->error);
	}

	if (FALSE == file_put_contents('bug53503.data', "1\n2\n3\n"))
		printf("[004] Failed to create CVS file\n");

	if (!$link->query("SELECT 1 FROM DUAL"))
		printf("[005] [%d] %s\n", $link->errno, $link->error);

	if (!$link->query("LOAD DATA LOCAL INFILE 'bug53503.data' INTO TABLE test")) {
		printf("[006] [%d] %s\n", $link->errno, $link->error);
		echo "bug";
	} else {
		echo "done";
	}
	$link->close();
?>
