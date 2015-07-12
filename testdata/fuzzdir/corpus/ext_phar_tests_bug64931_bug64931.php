<?php 

echo "Test\n";

@unlink(__DIR__."/bug64931.phar");
$phar = new Phar(__DIR__."/bug64931.phar");
$phar->addFile(__DIR__."/src/.pharignore", ".pharignore");
try {
	$phar->addFile(__DIR__."/src/.pharignore", ".phar/gotcha");
} catch (Exception $e) {
	echo "CAUGHT: ". $e->getMessage() ."\n";
}

try {
	$phar->addFromString(".phar", "gotcha");
} catch (Exception $e) {
	echo "CAUGHT: ". $e->getMessage() ."\n";
}

try {
	$phar->addFromString(".phar//", "gotcha");
} catch (Exception $e) {
	echo "CAUGHT: ". $e->getMessage() ."\n";
}

try {
	$phar->addFromString(".phar\\", "gotcha");
} catch (Exception $e) {
	echo "CAUGHT: ". $e->getMessage() ."\n";
}

try {
	$phar->addFromString(".phar\0", "gotcha");
} catch (Exception $e) {
	echo "CAUGHT: ". $e->getMessage() ."\n";
}

?>
===DONE===
