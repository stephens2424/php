<?php
include "include.inc";

$php = get_cgi_path();
reset_env_vars();

$filename = dirname(__FILE__)."/006.test.php";

$code = '
<?php

$test = "var";

class test {
	private $var;
}

echo test::$var;

?>
';

file_put_contents($filename, $code);

var_dump(`"$php" -n -l "$filename"`);
var_dump(`"$php" -n -l some.unknown`);

$code = '
<?php

class test 
	private $var;
}

?>
';

file_put_contents($filename, $code);

var_dump(`"$php" -n -l "$filename" 2>/dev/null`);

@unlink($filename);

echo "Done\n";
?>
