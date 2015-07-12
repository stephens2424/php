<?php

include "include.inc";

$php = get_cgi_path();
reset_env_vars();

var_dump(`$php -n -v`);

echo "Done\n";
?>
