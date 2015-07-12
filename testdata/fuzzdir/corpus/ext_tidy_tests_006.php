<?php 
	$a = tidy_parse_string("<HTML><asd asdf></HTML>");
	echo tidy_get_error_buffer($a);
	
?>
