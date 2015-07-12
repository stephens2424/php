<?php
$php = getenv("TEST_PHP_EXECUTABLE");
$socket_file = tempnam(sys_get_temp_dir(), pathinfo(__FILE__, PATHINFO_FILENAME) . '.sock');
$test_file = dirname(__FILE__) . '/' . pathinfo(__FILE__, PATHINFO_FILENAME) . '.inc';
if (file_exists($socket_file)) {
	unlink($socket_file);
}
$socket = stream_socket_server('unix://' . $socket_file);
var_dump($socket);
if (!$socket) {
	exit(1);
}
$desc = array(
	0 => $socket,
	1 => STDOUT,
	2 => STDERR,
);
$pipes = array();
$proc = proc_open("$php -n " . escapeshellarg($test_file), $desc, $pipes);
var_dump($proc);
if (!$proc) {
	exit(1);
}

$client_socket = stream_socket_client('unix://' . $socket_file);
var_dump($client_socket);
echo stream_get_contents($client_socket);
fclose($client_socket);

proc_terminate($proc);
proc_close($proc);
unlink($socket_file);
?>
