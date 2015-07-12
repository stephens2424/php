<?php

$php = getenv('TEST_PHP_EXECUTABLE');


echo `"$php" -n -m`;

echo "Done\n";
?>
