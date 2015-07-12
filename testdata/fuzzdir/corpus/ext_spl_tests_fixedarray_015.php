<?php

try {
	$a = new SplFixedArray('');
} catch (TypeError $iae) {
	echo "Ok - ".$iae->getMessage().PHP_EOL;
}

echo "Done\n";
?>
