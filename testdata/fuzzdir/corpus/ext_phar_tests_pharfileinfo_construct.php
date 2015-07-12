<?php
$fname = dirname(__FILE__) . '/' . basename(__FILE__, '.php') . '.phar';
$pname = 'phar://' . $fname;

try {
file_put_contents($fname, 'blah');
$a = new PharFileInfo($pname . '/oops');
} catch (Exception $e) {
echo $e->getMessage() . "\n";
unlink($fname);
}

try {
$a = new PharFileInfo(array());
} catch (TypeError $e) {
echo $e->getMessage() . "\n";
}

$a = new Phar($fname);
$a['a'] = 'hi';
$b = $a['a'];

try {
$a = new PharFileInfo($pname . '/oops/I/do/not/exist');
} catch (Exception $e) {
echo $e->getMessage() . "\n";
}

try {
$b->__construct('oops');
} catch (Exception $e) {
echo $e->getMessage() . "\n";
}

try {
$a = new PharFileInfo(__FILE__);
} catch (Exception $e) {
echo $e->getMessage() . "\n";
}
?>
===DONE===
