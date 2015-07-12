<?php

try {
	$a = new SplFixedArray('FOO');
} catch (TypeError $iae) {
	echo "Ok - ".$iae->getMessage().PHP_EOL;
}
?>
