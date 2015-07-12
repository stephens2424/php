<?php 
$key = ftok(dirname(__FILE__) . "/001.phpt", "p");
$q = msg_get_queue($key);
msg_send($q, 1, "hello") or print "FAIL\n";
$type = null;
if (msg_receive($q, 0, $type, 1024, $message)) {
	echo "TYPE: $type\n";
	echo "DATA: $message\n";
}
if (!msg_remove_queue($q)) {
	echo "BAD: queue removal failed\n";
}
?>
