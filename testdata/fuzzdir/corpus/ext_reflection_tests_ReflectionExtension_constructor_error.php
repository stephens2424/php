<?php
try {
	$obj = new ReflectionExtension();
} catch (TypeError $re) {
	echo "Ok - ".$re->getMessage().PHP_EOL;
}

try {
	$obj = new ReflectionExtension('foo', 'bar');
} catch (TypeError $re) {
	echo "Ok - ".$re->getMessage().PHP_EOL;
}

try {
	$obj = new ReflectionExtension([]);
} catch (TypeError $re) {
	echo "Ok - ".$re->getMessage().PHP_EOL;
}


?>
==DONE==
