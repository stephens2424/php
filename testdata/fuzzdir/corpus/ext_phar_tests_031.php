<?php

$pharconfig = 3;

require_once 'files/phar_oo_test.inc';

Phar::loadPhar($fname);

$pname = 'phar://' . $fname . '/a.php';

var_dump(file_get_contents($pname));
require $pname;

?>
===DONE===
