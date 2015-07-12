<?php
        $tidy = tidy_parse_file(dirname(__FILE__)."/016.html", dirname(__FILE__)."/016.tcfg");
    	tidy_clean_repair($tidy);
        echo tidy_get_output($tidy);
?>
