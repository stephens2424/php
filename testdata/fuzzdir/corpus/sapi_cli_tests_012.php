<?php

$php = getenv('TEST_PHP_EXECUTABLE');

var_dump(`"$php" -n -F some.php -F some.php`);
var_dump(`"$php" -n -F some.php -R some.php`);
var_dump(`"$php" -n -R some.php -F some.php`);
var_dump(`"$php" -n -R some.php -R some.php`);
var_dump(`"$php" -n -f some.php -f some.php`);
var_dump(`"$php" -n -B '' -B ''`);
var_dump(`"$php" -n -E '' -E ''`);
var_dump(`"$php" -n -r '' -r ''`);

echo "Done\n";
?>
