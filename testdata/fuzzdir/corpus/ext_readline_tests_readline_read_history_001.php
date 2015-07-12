<?php

$name = tempnam('/tmp', 'readline.tmp');

readline_add_history("foo\n");

var_dump(readline_write_history($name));

var_dump(readline_clear_history());

var_dump(readline_read_history($name));

var_dump(readline_list_history());

unlink($name);

?>
