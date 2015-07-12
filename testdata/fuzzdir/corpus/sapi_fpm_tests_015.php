<?php

include "include.inc";

$logfile = dirname(__FILE__).'/php-fpm.log.tmp';
$port1 = 9000+PHP_INT_SIZE;
$port2 = 9001+PHP_INT_SIZE;

$cfg = <<<EOT
[global]
error_log = $logfile
log_level = notice
[pool1]
listen = 127.0.0.1:$port1
listen.allowed_clients=127.0.0.1
user = foo
pm = dynamic
pm.max_children = 5
pm.min_spare_servers = 1
pm.max_spare_servers = 3
catch_workers_output = yes
[pool2]
listen = 127.0.0.1:$port2
listen.allowed_clients=xxx
pm = dynamic
pm.max_children = 5
pm.start_servers = 1
pm.min_spare_servers = 1
pm.max_spare_servers = 3
catch_workers_output = yes
EOT;

$fpm = run_fpm($cfg, $tail);
if (is_resource($fpm)) {
    $i = 0;
	while (($i++ < 30) && !($fp = @fsockopen('127.0.0.1', $port1))) {
		usleep(10000);
	}
	if ($fp) {
		echo "Started\n";
		fclose($fp);
	}
	for ($i=0 ; $i<10 ; $i++) {
		try {
			run_request('127.0.0.1', $port1);
		} catch (Exception $e) {
			echo "Error 1\n";
		}
	}
	try {
		run_request('127.0.0.1', $port2);
	} catch (Exception $e) {
		echo "Error 2\n";
	}
	proc_terminate($fpm);
	if (!feof($tail)) {
		echo stream_get_contents($tail);
	}
	fclose($tail);
	proc_close($fpm);
}

?>
Done
