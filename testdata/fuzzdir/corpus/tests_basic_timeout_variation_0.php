<?php

include dirname(__FILE__) . DIRECTORY_SEPARATOR . "timeout_config.inc";

set_time_limit($t);

while (1) { 
	busy_wait(1);
}

?>
never reached here
