<?php
	require('table.inc');

	$stmt = mysqli_stmt_init($link);
	mysqli_stmt_prepare($stmt, 'SELECT * FROM test');
	if (1 !== ($tmp = mysqli_stmt_attr_get($stmt, MYSQLI_STMT_ATTR_PREFETCH_ROWS))) {
		printf("[001] Expecting int/1, got %s/%s for attribute %s/%s\n",
			gettype($tmp), $tmp, $k, $attr);
	}
	$stmt->close();
	mysqli_close($link);
	print "done!";
?>
