<?php

$php = getenv('TEST_PHP_EXECUTABLE');

$filename = dirname(__FILE__)."/014.test.php";
$code = '
<?php
$test = "var"; //var
/* test class */
class test {
	private $var = array();

	public static function foo(Test $arg) {
		echo "hello";
		var_dump($this);
	}
}

$o = new test;
?>
';

file_put_contents($filename, $code);

var_dump(`"$php" -n -s $filename`);
var_dump(`"$php" -n -s unknown`);

@unlink($filename);

echo "Done\n";
?>
