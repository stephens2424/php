<?php
include "include.inc";

$php = get_cgi_path();
reset_env_vars();

var_dump(`"$php" -n -f some.php -f some.php`);
var_dump(`"$php" -n -s -w -l`);

?>
===DONE===
