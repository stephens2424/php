<?php

include "include.inc";

$php = get_cgi_path();
reset_env_vars();

$filename = dirname(__FILE__)."/008.test.php";
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

var_dump(`"$php" -n -s "$filename"`);
var_dump(`"$php" -n -s "unknown"`);

@unlink($filename);

echo "Done\n";
?>
