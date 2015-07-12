<?php

include dirname(__FILE__) . "/proc_open_pipes.inc";

for ($i = 3; $i<= 5; $i++) {
	$spec[$i] = array('pipe', 'w');
}

$php = getenv("TEST_PHP_EXECUTABLE");
$callee = create_sleep_script();

$spec[$i] = array('pi');
proc_open("$php $callee", $spec, $pipes);

$spec[$i] = 1;
proc_open("$php $callee", $spec, $pipes);

$spec[$i] = array('pipe', "test");
proc_open("$php $callee", $spec, $pipes);
var_dump($pipes);

$spec[$i] = array('file', "test", "z");
proc_open("$php $callee", $spec, $pipes);
var_dump($pipes);

echo "END\n";
?>
