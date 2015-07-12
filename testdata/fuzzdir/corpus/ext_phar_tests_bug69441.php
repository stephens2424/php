<?php
$fname = dirname(__FILE__) . '/bug69441.phar';
try {
$r = new Phar($fname, 0);
} catch(UnexpectedValueException $e) {
	echo $e;
}
?>

==DONE==
