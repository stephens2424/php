<?php
require_once('connect.inc');

$mysqli = mysqli_init();
$mysqli->real_connect($host, $user, $passwd, $db, $port, $socket);
if (mysqli_connect_errno()) {
	printf("Connect failed: %s\n", mysqli_connect_error());
	exit();
}

$mysqli->query("DROP PROCEDURE IF EXISTS p1") or die($mysqli->error);
$mysqli->query("CREATE PROCEDURE p1() BEGIN SELECT 23; SELECT 42; END") or die($mysqli->error);

if ($mysqli->multi_query("CALL p1();"))
{
	do
	{
		if ($objResult = $mysqli->store_result()) {
			while ($row = $objResult->fetch_assoc()) {
				print_r($row);
			}
			$objResult->close();
			if ($mysqli->more_results()) {
				print "----- next result -----------\n";
			}
		} else {
			print "no results found\n";
		}
	} while ($mysqli->more_results() && $mysqli->next_result());
} else {
	print $mysqli->error;
}

$mysqli->query("DROP PROCEDURE p1") or die($mysqli->error);
$mysqli->close();
print "done!";
?>
