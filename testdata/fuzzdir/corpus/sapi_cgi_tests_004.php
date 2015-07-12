<?php

include "include.inc";

$php = get_cgi_path();
reset_env_vars();

$filename = dirname(__FILE__).'/004.test.php';
$code ='
<?php

class test { 
	private $pri; 
}

var_dump(test::$pri);
?>
';

file_put_contents($filename, $code);

var_dump(`$php -n -f "$filename" 2>/dev/null`);
var_dump(`$php -n -f "wrong"`);

@unlink($filename);

echo "Done\n";
?>
