<?php

$a = 'a\\';
$b = 'b -c d\\';
var_dump( $a, escapeshellarg($a) );
var_dump( $b, escapeshellarg($b) );

$helper_script = <<<SCRIPT
<?php

print( "--- ARG INFO ---\n" );
var_dump( \$argv );

SCRIPT;

$script = dirname(__FILE__) . DIRECTORY_SEPARATOR . "arginfo.php";
file_put_contents($script, $helper_script);

$cmd =  PHP_BINARY . " " . $script . " "  . escapeshellarg($a) . " " . escapeshellarg($b);

system($cmd);

unlink($script);
?>
