<?php

$fname = dirname(__FILE__) . DIRECTORY_SEPARATOR . "bug69320.txt";
file_put_contents($fname, "foo");
var_dump(finfo_file(finfo_open(FILEINFO_MIME_TYPE), $fname));

?>
