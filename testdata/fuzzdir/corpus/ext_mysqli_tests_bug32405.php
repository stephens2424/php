<?php
	require_once("connect.inc");

	/*** test mysqli_connect 127.0.0.1 ***/
	$link = my_mysqli_connect($host, $user, $passwd, $db, $port, $socket);
	mysqli_select_db($link, "test");
	mysqli_query($link, "SET sql_mode=''");

	/* two fields are needed. the problem does not occur with 1 field only selected. */
	$link->query("CREATE TABLE test_users(user_id int(10) unsigned NOT NULL auto_increment, login varchar(50) default '', PRIMARY KEY (user_id))");
	$link->query('INSERT INTO test_users VALUES (NULL, "user1"), (NULL, "user2"), (NULL, "user3"), (NULL, "user4")');


	if ($stmt = $link->prepare("SELECT SQL_NO_CACHE user_id, login FROM test_users")) {
			$stmt->execute();
				$stmt->bind_result($col1, $col2);
				while ($stmt->fetch()) {
					var_dump($col1, $col2);
			}
			$stmt->close();
	}

	mysqli_query($link,"DROP TABLE test_users");
	mysqli_close($link);
?>
