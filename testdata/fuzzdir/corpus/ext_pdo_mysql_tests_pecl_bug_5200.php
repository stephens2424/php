<?php
require dirname(__FILE__) . '/../../../ext/pdo/tests/pdo_test.inc';
$db = PDOTest::test_factory(dirname(__FILE__). '/common.phpt');

$db->exec("CREATE TABLE test (bar INT NOT NULL, phase enum('please_select', 'I', 'II', 'IIa', 'IIb', 'III', 'IV'))");

foreach ($db->query('DESCRIBE test phase')->fetchAll(PDO::FETCH_ASSOC) as $row) {
	print_r($row);
}
?>
