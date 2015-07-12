<?php

$php = getenv('TEST_PHP_EXECUTABLE');

$filename = dirname(__FILE__).'/008.test.php';
$code ='
<?php

class test { 
	private $pri; 
}

var_dump(test::$pri);
?>
';

file_put_contents($filename, $code);

var_dump(`$php -n -f "$filename"`);
var_dump(`$php -n -f "wrong"`);

@unlink($filename);

echo "Done\n";
?>
