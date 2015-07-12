<?php
	try {
	    $section = new DOMCDataSection();
	} catch (TypeError $e) {
	    echo $e->getMessage();
	}
?>
