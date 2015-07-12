<?php

try {
    $nx = new Phar();
	$nx->getLinkTarget();
} catch (TypeError $e) {
	echo $e->getMessage(), "\n";
}

?>
