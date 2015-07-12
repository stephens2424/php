<?php

include dirname(__FILE__) . "/proc_open_pipes.inc";

for ($i = 3; $i<= 30; $i++) {
	$spec[$i] = array('pipe', 'w');
}

$php = getenv("TEST_PHP_EXECUTABLE");
$callee = create_sleep_script();
proc_open("$php $callee", $spec, $pipes);

var_dump(count($spec));
var_dump($pipes);

?>
