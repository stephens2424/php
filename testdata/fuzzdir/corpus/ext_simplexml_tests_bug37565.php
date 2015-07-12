<?php

function my_error_handler($errno, $errstr, $errfile, $errline) {
	    echo "Error: $errstr\n";
}

set_error_handler('my_error_handler');

class Setting extends ReflectionObject
{
}

try {
	Reflection::export(simplexml_load_string('<test/>', 'Setting'));
} catch (Error $e) {
	my_error_handler($e->getCode(), $e->getMessage(), $e->getFile(), $e->getLine());
}

try {
	Reflection::export(simplexml_load_file('data:,<test/>', 'Setting'));
} catch (Error $e) {
	my_error_handler($e->getCode(), $e->getMessage(), $e->getFile(), $e->getLine());
}

?>
===DONE===
