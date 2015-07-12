<?php 

	$a = tidy_parse_string("<HTML></HTML>");
	tidy_clean_repair($a);
	echo tidy_get_output($a);

?>
