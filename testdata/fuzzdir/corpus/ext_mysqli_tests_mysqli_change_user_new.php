<?php
	require_once("connect.inc");

	$tmp	= NULL;
	$link	= NULL;

	if (!$link = my_mysqli_connect($host, $user, $passwd, $db, $port, $socket))
		printf("[001] Cannot connect to the server using host=%s, user=%s, passwd=***, dbname=%s, port=%s, socket=%s\n",
			$host, $user, $db, $port, $socket);

	/* Pre 5.6: link remains useable */
	if (false !== ($tmp = @mysqli_change_user($link, $user . '_unknown_really', $passwd . 'non_empty', $db)))
		printf("[002] Expecting false, got %s/%s\n", gettype($tmp), $tmp);

	if (!$res = mysqli_query($link, 'SELECT 1 AS _one'))
		printf("[003] [%d] %s\n", mysqli_errno($link), mysqli_error($link));

	var_dump($res->fetch_assoc());

	print "done!";
?>
