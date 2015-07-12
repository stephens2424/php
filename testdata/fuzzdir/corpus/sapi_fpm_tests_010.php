<?php

include "include.inc";

$logfile = dirname(__FILE__).'/php-fpm.log.tmp';
$port = 9000+PHP_INT_SIZE;

$cfg = <<<EOT
[global]
error_log = $logfile
[unconfined]
listen = 127.0.0.1:$port
pm.status_path = /status
pm = static 
pm.max_children = 1
EOT;

$fpm = run_fpm($cfg, $tail);
if (is_resource($fpm)) {
    fpm_display_log($tail, 2);
    try {
		echo run_request('127.0.0.1', $port, '/status');

		$html = run_request('127.0.0.1', $port, '/status', 'html');
		var_dump(strpos($html, 'text/html') && strpos($html, 'DOCTYPE') && strpos($html, 'PHP-FPM Status Page'));

		$json = run_request('127.0.0.1', $port, '/status', 'json');
		var_dump(strpos($json, 'application/json') && strpos($json, '"pool":"unconfined"'));

		$xml = run_request('127.0.0.1', $port, '/status', 'xml');
		var_dump(strpos($xml, 'text/xml') && strpos($xml, '<?xml'));

		echo "IPv4 ok\n";
	} catch (Exception $e) {
		echo "IPv4 error\n";
	}

	proc_terminate($fpm);
    stream_get_contents($tail);
    fclose($tail);
    proc_close($fpm);
}

?>
