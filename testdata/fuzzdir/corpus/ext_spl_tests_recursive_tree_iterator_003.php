<?php
try {
	new RecursiveTreeIterator(new ArrayIterator(array()));
} catch (TypeError $e) {
    echo $e->getMessage(), "\n";
}
?>
===DONE===
