<?php

include "include.inc";

$php = get_cgi_path();
reset_env_vars();

var_dump(`$php -n -a -f 'wrong'`);
var_dump(`$php -n -f 'wrong' -a`);

echo "Done\n";
?>
