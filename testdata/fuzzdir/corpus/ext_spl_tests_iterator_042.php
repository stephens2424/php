<?php

function test_error_handler($errno, $msg, $filename, $linenum)
{
	echo "Error $msg in $filename on line $linenum\n";
	return true;
}

set_error_handler('test_error_handler');

$it = new AppendIterator;

try {
	$it->append(array());
} catch (Error $e) {
	test_error_handler($e->getCode(), $e->getMessage(), $e->getFile(), $e->getLine());
}
$it->append(new ArrayIterator(array(1)));
$it->append(new ArrayIterator(array(21, 22)));

var_dump($it->getArrayIterator());

$it->append(new ArrayIterator(array(31, 32, 33)));

var_dump($it->getArrayIterator());

$idx = 0;

foreach($it as $k => $v)
{
	echo '===' . $idx++ . "===\n";
	var_dump($it->getIteratorIndex());
	var_dump($k);
	var_dump($v);
}

?>
===DONE===
