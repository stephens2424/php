<?php 
   	$a = tidy_parse_string("<HTML><asd asdf></HTML>");
	echo $a->errorBuffer;
?>
