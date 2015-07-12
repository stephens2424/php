<?php

$php = getenv('TEST_PHP_EXECUTABLE');

var_dump(`"$php" -n --rc unknown`);
var_dump(`"$php" -n --rc stdclass`);
var_dump(`"$php" -n --rc exception`);

echo "Done\n";
?>
