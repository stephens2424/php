<?php

$php = getenv('TEST_PHP_EXECUTABLE');


echo `"$php" -n --ri this_extension_does_not_exist_568537753423`;
echo `"$php" -n --ri standard`;

echo "\nDone\n";
?>
