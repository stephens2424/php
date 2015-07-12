<?php

$php = getenv('TEST_PHP_EXECUTABLE');

$filename = dirname(__FILE__).'/007.test.php';
$code ='
<?php
/* some test script */

class test { /* {{{ */
	public $var = "test"; //test var
#perl style comment 
	private $pri; /* private attr */

	function foo(/* void */) {
	}
}
/* }}} */

?>
';

file_put_contents($filename, $code);

var_dump(`$php -n -w "$filename"`);
var_dump(`$php -n -w "wrong"`);
var_dump(`echo "<?php /* comment */ class test {\n // comment \n function foo() {} } ?>" | $php -n -w`);

@unlink($filename);

echo "Done\n";
?>
