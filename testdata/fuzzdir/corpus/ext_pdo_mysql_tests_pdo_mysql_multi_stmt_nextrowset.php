<?php
	require_once(dirname(__FILE__) . DIRECTORY_SEPARATOR . 'mysql_pdo_test.inc');
	$db = MySQLPDOTest::factory();
	$db->setAttribute(PDO::ATTR_STRINGIFY_FETCHES, true);

	MySQLPDOTest::createTestTable($db);

	function test_proc($db) {

		$db->exec('DROP PROCEDURE IF EXISTS p');
		$db->exec('CREATE PROCEDURE p() BEGIN SELECT id FROM test ORDER BY id ASC LIMIT 3; SELECT id, label FROM test WHERE id < 4 ORDER BY id DESC LIMIT 3; END;');
		$stmt = $db->query('CALL p()');
		do {
			var_dump($stmt->fetchAll(PDO::FETCH_ASSOC));
		} while ($stmt->nextRowSet());
		var_dump($stmt->nextRowSet());

	}

	try {

		// Using native PS for proc, since emulated fails.
		printf("Native PS...\n");
		foreach (array(false, true) as $multi) {
			$value = $multi ? 'true' : 'false';
			echo "\nTesting with PDO::MYSQL_ATTR_MULTI_STATEMENTS set to {$value}\n";
			$dsn = MySQLPDOTest::getDSN();
			$user = PDO_MYSQL_TEST_USER;
			$pass = PDO_MYSQL_TEST_PASS;
			$db = new PDO($dsn, $user, $pass, array(PDO::MYSQL_ATTR_MULTI_STATEMENTS => $multi));
			$db->setAttribute(PDO::ATTR_STRINGIFY_FETCHES, true);
			$db->setAttribute(PDO::MYSQL_ATTR_USE_BUFFERED_QUERY, 1);
			$db->setAttribute(PDO::ATTR_EMULATE_PREPARES, 0);
			test_proc($db);

			$db = new PDO($dsn, $user, $pass, array(PDO::MYSQL_ATTR_MULTI_STATEMENTS => $multi));
			$db->setAttribute(PDO::ATTR_STRINGIFY_FETCHES, true);
			$db->setAttribute(PDO::MYSQL_ATTR_USE_BUFFERED_QUERY, 0);
			$db->setAttribute(PDO::ATTR_EMULATE_PREPARES, 0);

			test_proc($db);

			// Switch back to emulated prepares to verify multi statement attribute.
			$db->setAttribute(PDO::ATTR_EMULATE_PREPARES, 1);
			// This will fail when $multi is false.
			$stmt = $db->query("SELECT * FROM test; INSERT INTO test (id, label) VALUES (99, 'x')");
			if ($stmt !== false) {
				$stmt->closeCursor();
			}
			$info = $db->errorInfo();
			var_dump($info[0]);
		}
		@$db->exec('DROP PROCEDURE IF EXISTS p');

	} catch (PDOException $e) {
		printf("[001] %s [%s] %s\n",
			$e->getMessage(), $db->errorCode(), implode(' ', $db->errorInfo()));
	}

	print "done!";
?>
