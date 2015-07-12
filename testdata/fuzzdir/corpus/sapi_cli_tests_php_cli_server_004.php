<?php
include "php_cli_server.inc";
php_cli_server_start('foreach($_SERVER as $k=>$v) { if (!strncmp($k, "HTTP", 4)) var_dump( $k . ":" . $v); }');

list($host, $port) = explode(':', PHP_CLI_SERVER_ADDRESS);
$port = intval($port)?:80;

$fp = fsockopen($host, $port, $errno, $errstr, 0.5);
if (!$fp) {
  die("connect failed");
}

if(fwrite($fp, <<<HEADER
GET / HTTP/1.1
Host:{$host}
User-Agent:dummy
Custom:foo
Referer:http://www.php.net/


HEADER
)) {
	while (!feof($fp)) {
		echo fgets($fp);
	}
}

?>
