<?php

        function dump_nodes(tidyNode $node) {

            var_dump($node->hasChildren());
            if($node->hasChildren()) {

                foreach($node->child as $c) {

                    var_dump($c);

                    if($c->hasChildren()) {

                        dump_nodes($c);

                    }
                }

            }

        }

    	$a = tidy_parse_string("<HTML><BODY BGCOLOR=#FFFFFF ALINK=#000000><B>Hi</B><I>Bye<U>Test</U></I></BODY></HTML>", array('newline' => 'LF'));
        $html = $a->html();
        dump_nodes($html);
            
?>
