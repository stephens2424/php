<?php

$php = getenv('TEST_PHP_EXECUTABLE');

var_dump(`$php -n --re unknown`);
var_dump(`$php -n --re ""`);
var_dump(`$php -n --re pcre`);

echo "Done\n";
?>
