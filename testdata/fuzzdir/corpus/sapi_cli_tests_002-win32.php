<?php

$php = getenv('TEST_PHP_EXECUTABLE');

var_dump(`$php -n -r "var_dump('hello');"`);

echo "Done\n";
?>
