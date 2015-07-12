<?php

try {
	new SplFileObject('foo', array());
} catch (TypeError $e) {
	var_dump($e->getMessage());
}

?>
