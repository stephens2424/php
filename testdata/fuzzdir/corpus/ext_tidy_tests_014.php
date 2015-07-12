<?php
        $text = "<B>testing</I>";
    	$tidy = tidy_parse_string($text, array('show-body-only'=>true));
    	tidy_clean_repair($tidy);
    	echo tidy_get_output($tidy);
            
?>
