<?php
include "php_cli_server.inc";
php_cli_server_start('header(\'WWW-Authenticate: Digest realm="foo",qop="auth",nonce="XXXXX",opaque="'.md5("foo").'"\');');

list($host, $port) = explode(':', PHP_CLI_SERVER_ADDRESS);
$port = intval($port)?:80;

$fp = fsockopen($host, $port, $errno, $errstr, 0.5);
if (!$fp) {
  die("connect failed");
}

if(fwrite($fp, <<<HEADER
GET / HTTP/1.1
Host: {$host}
Authorization: Basic Zm9vOmJhcg==


HEADER
)) {
	while (!feof($fp)) {
		echo fgets($fp);
	}
}

?>
