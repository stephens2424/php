<?php

$php = getenv('TEST_PHP_EXECUTABLE');


echo `"$php" -n --version | grep built:`;
echo `echo "<?php print_r(\\\$argv);" | "$php" -n -- foo bar baz`, "\n";
echo `"$php" -n --version foo bar baz | grep built:`;
echo `"$php" -n --notexisting foo bar baz | grep Usage:`;

echo "Done\n";
?>
