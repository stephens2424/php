<?php

try {
	$a = new ReflectionFunction(array(1, 2, 3));
	echo "exception not thrown.".PHP_EOL;
} catch (TypeError $re) {
	echo "Ok - ".$re->getMessage().PHP_EOL;
}
try {
	$a = new ReflectionFunction('nonExistentFunction');
} catch (ReflectionException $e) {
	echo $e->getMessage().PHP_EOL;
}
try {
	$a = new ReflectionFunction();
} catch (TypeError $re) {
	echo "Ok - ".$re->getMessage().PHP_EOL;
}
try {
	$a = new ReflectionFunction(1, 2);
} catch (TypeError $re) {
	echo "Ok - ".$re->getMessage().PHP_EOL;
}
try {
	$a = new ReflectionFunction([]);
} catch (TypeError $re) {
	echo "Ok - ".$re->getMessage().PHP_EOL;
}

?>
