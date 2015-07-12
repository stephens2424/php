<?php

include dirname(__FILE__) . DIRECTORY_SEPARATOR . "timeout_config.inc";

set_time_limit($t);

function f()
{
	echo "call";
	busy_wait(4);
}

register_shutdown_function("f");
?>
shutdown happens after here
