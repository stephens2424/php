<?php
try {
	$p = new PharData(dirname(__FILE__) . '/files/trunc.tar');
} catch (Exception $e) {
	echo $e->getMessage() . "\n";
}

?>
===DONE===
