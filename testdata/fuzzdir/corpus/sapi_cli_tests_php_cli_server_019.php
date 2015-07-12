<?php
include "php_cli_server.inc";
php_cli_server_start(<<<'PHP'
header('Bar-Foo: Foo');
var_dump(getallheaders());
var_dump(apache_request_headers());
var_dump(apache_response_headers());
PHP
);

list($host, $port) = explode(':', PHP_CLI_SERVER_ADDRESS);
$port = intval($port)?:80;

$fp = fsockopen($host, $port, $errno, $errstr, 0.5);
if (!$fp) {
  die("connect failed");
}

if(fwrite($fp, <<<HEADER
GET / HTTP/1.1
Host: {$host}
Foo-Bar: Bar


HEADER
)) {
    while (!feof($fp)) {
        echo fgets($fp);
    }
}

fclose($fp);
?>
