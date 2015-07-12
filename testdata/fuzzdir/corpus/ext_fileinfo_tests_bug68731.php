<?php
	$buffer = file_get_contents(dirname(__FILE__) . '/68731.gif');
	$finfo = finfo_open(FILEINFO_MIME_TYPE);
	echo finfo_buffer($finfo, $buffer);
?>
