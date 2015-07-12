<?php
$dir = 'file://' . dirname(__FILE__) . '/testDir';
mkdir($dir);
var_dump(is_dir($dir));
rmdir($dir);
var_dump(is_dir($dir));
?>
