<?php

$php = getenv('TEST_PHP_EXECUTABLE');

$filename_txt = dirname(__FILE__)."/013.test.txt";
file_put_contents($filename_txt, "test\nfile\ncontents\n");

var_dump(`cat "$filename_txt" | "$php" -n -B 'var_dump("start");'`);
var_dump(`cat "$filename_txt" | "$php" -n -E 'var_dump("end");'`);
var_dump(`cat "$filename_txt" | "$php" -n -B 'var_dump("start");' -E 'var_dump("end");'`);

@unlink($filename_txt);

echo "Done\n";
?>
