<?php
$fname = dirname(__FILE__) . '/bug69453.tar.phar';
try {
$r = new Phar($fname, 0);
} catch(UnexpectedValueException $e) {
	echo $e;
}
?>

==DONE==
