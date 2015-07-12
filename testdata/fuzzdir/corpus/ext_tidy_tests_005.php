<?php 
	$a = tidy_parse_file(dirname(__FILE__)."/005.html");
	echo tidy_get_output($a);
	
?>
