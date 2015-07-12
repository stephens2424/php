<?php
try {
	$array = new SplFixedArray(new SplFixedArray(3));
} catch (TypeError $iae) {
	echo "Ok - ".$iae->getMessage().PHP_EOL;
}

?>
