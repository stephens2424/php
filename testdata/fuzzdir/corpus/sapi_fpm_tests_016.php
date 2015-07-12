<?php

include "include.inc";

$logfile = __DIR__.'/php-fpm.log.tmp';
$logdir  = __DIR__.'/conf.d';
$port = 9000+PHP_INT_SIZE;

// Main configuration
$cfg = <<<EOT
[global]
error_log = $logfile
log_level = notice
include = $logdir/*.conf
EOT;

// Splited configuration
@mkdir($logdir);
$i=$port;
$names = ['cccc', 'aaaa', 'eeee', 'dddd', 'bbbb'];
foreach($names as $name) {
	$poolcfg = <<<EOT
[$name]
listen = 127.0.0.1:$i
listen.allowed_clients=127.0.0.1
user = foo
pm = ondemand
pm.max_children = 5
EOT;
	file_put_contents("$logdir/$name.conf", $poolcfg);
	$i++;
}

// Test
$fpm = run_fpm($cfg, $tail);
if (is_resource($fpm)) {
    fpm_display_log($tail, count($names)+2);
	$i=$port;
	foreach($names as $name) {
		try {
			run_request('127.0.0.1', $i++);
			echo "OK $name\n";
		} catch (Exception $e) {
			echo "Error 1\n";
		}
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
