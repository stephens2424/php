<?php

$php = getenv('TEST_PHP_EXECUTABLE');

var_dump(`$php -n -a -r "echo hello;"`);
var_dump(`$php -n -r "echo hello;" -a`);

echo "Done\n";
?>
