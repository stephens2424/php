<?php
include "php_cli_server.inc";
php_cli_server_start('echo done, "\n";', null);

list($host, $port) = explode(':', PHP_CLI_SERVER_ADDRESS);
$port = intval($port)?:80;
$output = '';

// note: select() on Windows (& some other platforms) has historical issues with
//       timeouts less than 1000 millis(0.5). it may be better to increase these
//       timeouts to 1000 millis(1.0) (fsockopen eventually calls select()).
//       see articles like: http://support.microsoft.com/kb/257821
$fp = fsockopen($host, $port, $errno, $errstr, 0.5);
if (!$fp) {
  die("connect failed");
}

if(fwrite($fp, <<<HEADER
POST /index.php HTTP/1.1
Host: {$host}
Content-Type: multipart/form-data; boundary=---------123456789
Content-Length: 70

---------123456789
Content-Type: application/x-www-form-urlencoded
a=b
HEADER
)) {
	while (!feof($fp)) {
		$output .= fgets($fp);
	}
}

fclose($fp);

$fp = fsockopen($host, $port, $errno, $errstr, 0.5);
if(fwrite($fp, <<<HEADER
POST /main/no-exists.php HTTP/1.1
Host: {$host}
Content-Type: multipart/form-data; boundary=---------123456789
Content-Length: 70

---------123456789
Content-Type: application/x-www-form-urlencoded
a=b
HEADER
)) {
	while (!feof($fp)) {
		$output .= fgets($fp);
	}
}

echo preg_replace("/<style>(.*?)<\/style>/s", "<style>AAA</style>", $output), "\n";
fclose($fp);

?>
