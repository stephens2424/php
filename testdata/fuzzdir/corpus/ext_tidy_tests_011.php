<?php 
    	$a = tidy_parse_string("<HTML><BODY BGCOLOR=#FFFFFF ALINK=#000000></BODY></HTML>");
        $body = $a->body();
        var_dump($body->attribute);
        foreach($body->attribute as $key=>$val) {
            echo "Attrib '$key': $val\n";
        }
?>
