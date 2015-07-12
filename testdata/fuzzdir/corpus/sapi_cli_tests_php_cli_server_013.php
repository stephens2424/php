<?php
include "php_cli_server.inc";
php_cli_server_start(NULL, NULL);

list($host, $port) = explode(':', PHP_CLI_SERVER_ADDRESS);
$port = intval($port)?:80;
$output = '';

$fp = fsockopen($host, $port, $errno, $errstr, 0.5);
if (!$fp) {
  die("connect failed");
}


if(fwrite($fp, <<<HEADER
POST / HTTP/1.1
Host: {$host}
Content-Type: application/x-www-form-urlencoded
Content-Length: 3

a=b
HEADER
)) {
	while (!feof($fp)) {
		$output .= fgets($fp);
	}
}

echo preg_replace("/<style>(.*?)<\/style>/s", "<style>AAA</style>", $output), "\n";
fclose($fp);


$output = '';
$fp = fsockopen($host, $port, $errno, $errstr, 0.5);
if (!$fp) {
  die("connect failed");
}

if(fwrite($fp, <<<HEADER
GET /main/style.css HTTP/1.1
Host: {$host}


HEADER
)) {
	while (!feof($fp)) {
		$output .= fgets($fp);
	}
}

echo preg_replace("/<style>(.*?)<\/style>/s", "<style>AAA</style>", $output), "\n";
fclose($fp);

$output = '';
$fp = fsockopen($host, $port, $errno, $errstr, 0.5);
if (!$fp) {
  die("connect failed");
}

if(fwrite($fp, <<<HEADER
HEAD /main/foo/bar HTTP/1.1
Host: {$host}


HEADER
)) {
	while (!feof($fp)) {
		$output .= fgets($fp);
	}
}

echo preg_replace("/<style>(.*?)<\/style>/s", "<style>AAA</style>", $output), "\n";
fclose($fp);
?>
