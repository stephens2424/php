<?php

function foo(): stdClass {
	$a = new stdClass;
	$b = [];
	return [$a, $b];
}

try {
	foo();
} catch (Error $e) {
	print $e->getMessage();
}

?>
