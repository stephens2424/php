<?php

class DbConnection {
	public function connect() {
		require_once("connect.inc");

		$link = my_mysqli_connect($host, $user, $passwd, $db, $port, $socket);
		var_dump($link);

		$link = mysqli_init();
		/* @ is to suppress 'Property access is not allowed yet' */
		@var_dump($link);

		$mysql = new my_mysqli($host, $user, $passwd, $db, $port, $socket);
		$mysql->query("DROP TABLE IF EXISTS test_warnings");
		$mysql->query("CREATE TABLE test_warnings (a int not null)");
		$mysql->query("SET sql_mode=''");
		$mysql->query("INSERT INTO test_warnings VALUES (1),(2),(NULL)");

		$warning = $mysql->get_warnings();
		if (!$warning)
			printf("[001] No warning!\n");

		if ($warning->errno == 1048 || $warning->errno == 1253) {
			/* 1048 - Column 'a' cannot be null, 1263 - Data truncated; NULL supplied to NOT NULL column 'a' at row */
			if ("HY000" != $warning->sqlstate)
				printf("[003] Wrong sql state code: %s\n", $warning->sqlstate);

			if ("" == $warning->message)
				printf("[004] Message string must not be empty\n");


		} else {
			printf("[002] Empty error message!\n");
			var_dump($warning);
		}
	}
}

$db = new DbConnection();
$db->connect();

echo "Done\n";
?>
