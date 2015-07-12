<?php

include "include.inc";

$php = get_cgi_path();
reset_env_vars();

$f = tempnam(sys_get_temp_dir(), 'cgitest');

function test($script) {
	file_put_contents($GLOBALS['f'], $script);
	$cmd = escapeshellcmd($GLOBALS['php']);
	$cmd .= ' -n -dreport_zend_debug=0 -dhtml_errors=0 ' . escapeshellarg($GLOBALS['f']);
	echo "----------\n";
	echo rtrim($script) . "\n";
	echo "----------\n";
	passthru($cmd);
}

test('<?php ?>');
test('<?php header_remove(); ?>');
test('<?php header_remove("X-Foo"); ?>');
test('<?php
header("X-Foo: Bar");
?>');
test('<?php
header("X-Foo: Bar");
header("X-Bar: Baz");
header_remove("X-Foo");
?>');
test('<?php
header("X-Foo: Bar");
header_remove("X-Foo: Bar");
?>');
test('<?php
header("X-Foo: Bar");
header_remove("X-Foo:");
?>');
test('<?php
header("X-Foo: Bar");
header_remove();
?>');
test('<?php
header_remove("");
?>');
test('<?php
header_remove(":");
?>');
test('<?php
header("X-Foo: Bar");
echo "flush\n";
flush();
header_remove("X-Foo");
?>');

@unlink($f);
?>
