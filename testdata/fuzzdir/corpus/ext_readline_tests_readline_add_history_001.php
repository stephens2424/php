<?php

var_dump(readline_add_history('foo'));
var_dump(readline_list_history());
var_dump(readline_add_history(NULL));
var_dump(readline_list_history());
var_dump(readline_clear_history());
var_dump(readline_add_history());

?>
