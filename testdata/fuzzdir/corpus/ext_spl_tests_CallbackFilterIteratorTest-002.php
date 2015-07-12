<?php

set_error_handler(function($errno, $errstr){
	echo $errstr . "\n";
	return true;
});

try {
	new CallbackFilterIterator();
} catch (TypeError $e) {
	echo $e->getMessage() . "\n";
}

try {
	new CallbackFilterIterator(null);
} catch (TypeError $e) {
	echo $e->getMessage() . "\n";
}

try {
	new CallbackFilterIterator(new ArrayIterator(array()), null);
} catch (TypeError $e) {
	echo $e->getMessage() . "\n";
}

try {
	new CallbackFilterIterator(new ArrayIterator(array()), array());
} catch (TypeError $e) {
	echo $e->getMessage() . "\n";
}

$it = new CallbackFilterIterator(new ArrayIterator(array(1)), function() {
	throw new Exception("some message");
});
try {
	foreach($it as $e);
} catch(Exception $e) {
	echo $e->getMessage() . "\n";
}

