<?php

var_dump(readline_callback_handler_remove());
var_dump(readline_callback_handler_install('testing: ', 'foo'));

function foo() { }
var_dump(readline_callback_handler_install('testing: ', 'foo'));
var_dump(readline_callback_handler_remove());

?>
