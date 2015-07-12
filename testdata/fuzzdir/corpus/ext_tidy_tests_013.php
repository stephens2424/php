<?php
        $tidy = new tidy(dirname(__FILE__)."/013.html", array("show-body-only"=>true));
        $tidy->cleanRepair();
        echo $tidy;

?>
