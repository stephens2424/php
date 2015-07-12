<?php

include "include.inc";

$php = get_cgi_path();
reset_env_vars();

$f = tempnam(sys_get_temp_dir(), 'cgitest');

putenv("TRANSLATED_PATH=".$f."/x");
putenv("SCRIPT_FILENAME=".$f."/x");
file_put_contents($f, '<?php
header("HTTP/1.1 403 Forbidden");
header("Status: 403 Also Forbidden");
?>');

echo (`$php -n $f`);

file_put_contents($f, '<?php
header("HTTP/1.1 403 Forbidden");
?>');

echo (`$php -n $f`);

file_put_contents($f, '<?php
header("Status: 403 Also Forbidden");
?>');

echo (`$php -n $f`);

echo "Done\n";

@unlink($f);
?>
