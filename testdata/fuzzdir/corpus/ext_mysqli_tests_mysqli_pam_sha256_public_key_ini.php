<?php
	require_once("connect.inc");


	$link = new mysqli($host, 'shatest', 'shatest', $db, $port, $socket);
	if ($link->connect_errno) {
		printf("[001] [%d] %s\n", $link->connect_errno, $link->connect_error);
	} else {
		if (!$res = $link->query("SELECT id FROM test WHERE id = 1"))
			printf("[002] [%d] %s\n", $link->errno, $link->error);

		if (!$row = mysqli_fetch_assoc($res)) {
			printf("[003] [%d] %s\n", $link->errno, $link->error);
		}

		if ($row['id'] != 1) {
			printf("[004] Expecting 1 got %s/'%s'", gettype($row['id']), $row['id']);
		}
	}
	print "done!";
?>
