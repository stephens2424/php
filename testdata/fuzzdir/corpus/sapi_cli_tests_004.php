<?php

$php = getenv('TEST_PHP_EXECUTABLE');

var_dump(`$php -n --rf unknown`);
var_dump(`$php -n --rf echo`);
var_dump(`$php -n --rf phpinfo`);

echo "Done\n";
?>
