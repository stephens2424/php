<?php
	$test_file = dirname(__FILE__) . DIRECTORY_SEPARATOR . "bug68735.jpg";
	$f = new finfo;

	var_dump($f->file($test_file));

?>
===DONE===
