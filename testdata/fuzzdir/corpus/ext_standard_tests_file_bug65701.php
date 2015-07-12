<?php
$file_path = dirname(__FILE__) . "/bug65701/";

if (!is_dir($file_path)) {
	mkdir($file_path);
}

$src = $file_path . '/srcbug65701_file.txt';
$dst = tempnam($file_path, 'dstbug65701_file.txt');

file_put_contents($src, "Hello World");

copy($src, $dst);
var_dump(filesize($dst));
?>
