<?php
include "php_cli_server.inc";
php_cli_server_start('var_dump(count($_SERVER));', 'not-index.php');

list($host, $port) = explode(':', PHP_CLI_SERVER_ADDRESS);
$port = intval($port)?:80;

$fp = fsockopen($host, $port, $errno, $errstr, 0.5);
if (!$fp) {
  die("connect failed");
}

if(fwrite($fp, "GET www.example.com:80 HTTP/1.1\r\n\r\n")) {
    while (!feof($fp)) {
        echo fgets($fp);
    }
}

fclose($fp);
?>
