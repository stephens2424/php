<?php

$php = getenv('TEST_PHP_EXECUTABLE');


echo `"$php" -n -i`;

echo "\nDone\n";
?>
