<?php
        $tidy = tidy_parse_file(dirname(__FILE__)."/015.html", array('show-body-only'=>true));
    	tidy_clean_repair($tidy);
    	echo tidy_get_output($tidy);
            
?>
