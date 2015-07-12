<?php
	require_once("connect.inc");

	$tmp    = NULL;
	$link   = NULL;

	if (!is_null($tmp = @mysqli_sqlstate()))
		printf("[001] Expecting NULL, got %s/%s\n", gettype($tmp), $tmp);

	if (!is_null($tmp = @mysqli_sqlstate($link)))
		printf("[002] Expecting NULL, got %s/%s\n", gettype($tmp), $tmp);

	require('table.inc');

	var_dump(@mysqli_sqlstate($link, "foo"));

	var_dump(mysqli_sqlstate($link));
	mysqli_query($link, "SELECT unknown_column FROM test");
	var_dump(mysqli_sqlstate($link));
	mysqli_free_result(mysqli_query($link, "SELECT id FROM test"));
	var_dump(mysqli_sqlstate($link));

	mysqli_close($link);

	var_dump(mysqli_sqlstate($link));

	print "done!";
?>
