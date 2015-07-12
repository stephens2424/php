<?php
	require_once("connect.inc");

	/*** test mysqli_connect 127.0.0.1 ***/
	$link = my_mysqli_connect($host, $user, $passwd, $db, $port, $socket);

	mysqli_select_db($link, $db);

	mysqli_query($link, "DROP TABLE IF EXISTS test_affected");
	mysqli_query($link, "CREATE TABLE test_affected (foo int, bar varchar(10) character set latin1) ENGINE=" . $engine);

	mysqli_query($link, "INSERT INTO test_affected VALUES (1, 'Zak'),(2, 'Greant')");

	$stmt = mysqli_prepare($link, "SELECT * FROM test_affected");
	mysqli_stmt_execute($stmt);
	$result = mysqli_stmt_result_metadata($stmt);

	echo "\n=== fetch_fields ===\n";
	var_dump(mysqli_fetch_fields($result));

	echo "\n=== fetch_field_direct ===\n";
	var_dump(mysqli_fetch_field_direct($result, 0));
	var_dump(mysqli_fetch_field_direct($result, 1));

	echo "\n=== fetch_field ===\n";
	while ($field = mysqli_fetch_field($result)) {
		var_dump($field);
	}

	print_r(mysqli_fetch_lengths($result));

	mysqli_free_result($result);


	mysqli_stmt_close($stmt);
	mysqli_query($link, "DROP TABLE IF EXISTS test_affected");
	mysqli_close($link);
	print "done!";
?>
