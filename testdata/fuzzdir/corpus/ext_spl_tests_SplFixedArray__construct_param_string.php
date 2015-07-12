<?php
try {
	$array = new SplFixedArray( "string" );
} catch (TypeError $iae) {
	echo "Ok - ".$iae->getMessage().PHP_EOL;
}


?>
