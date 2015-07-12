<?php

include dirname(__FILE__) . DIRECTORY_SEPARATOR . "timeout_config.inc";

set_time_limit($t);

function f($t) { 
	echo "call";
	busy_wait($t-1);
	throw new Exception("exception before timeout");
}

f($t);
?>
never reached here
