<?php

$pharconfig = 0;

require_once 'files/phar_oo_test.inc';

$phar = new Phar($fname);
$phar->setInfoClass('SplFileObject');

$phar['hi/f.php'] = 'hi';
var_dump(isset($phar['hi']));
var_dump(isset($phar['hi/f.php']));
echo $phar['hi/f.php'];
echo "\n";

?>
===DONE===
