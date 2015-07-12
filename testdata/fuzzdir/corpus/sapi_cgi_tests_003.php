<?php

include "include.inc";

$php = get_cgi_path();
reset_env_vars();

$filename = dirname(__FILE__).'/003.test.php';
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
