<?php

try {
	$array = new SplFixedArray( array("string", 1) );
} catch (TypeError $iae) {
	echo "Ok - ".$iae->getMessage().PHP_EOL;
}

?>
