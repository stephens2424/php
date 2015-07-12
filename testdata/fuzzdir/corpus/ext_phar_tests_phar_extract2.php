<?php

$fname = dirname(__FILE__) . '/tempmanifest2.phar.php';
$pname = 'phar://' . $fname;

$phar = new Phar($fname);
$phar->setDefaultStub();
$phar->setAlias('fred');
$phar['file1.txt'] = 'hi';
$phar['file2.txt'] = 'hi2';
$phar['subdir/ectory/file.txt'] = 'hi3';
$phar->mount($pname . '/mount2', __FILE__);
$phar->addEmptyDir('one/level');

$phar->extractTo(dirname(__FILE__) . '/extract2', 'mount2');
$phar->extractTo(dirname(__FILE__) . '/extract2');
$out = array();

foreach (new RecursiveIteratorIterator(new RecursiveDirectoryIterator(dirname(__FILE__) . '/extract2', 0x00003000), RecursiveIteratorIterator::CHILD_FIRST) as $path => $file) {
	$extracted[] = $path;
}

sort($extracted);

foreach ($extracted as $out) {
	echo "$out\n";
}

?>
===DONE===
