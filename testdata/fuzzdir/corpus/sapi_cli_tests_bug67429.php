<?php
include "php_cli_server.inc";

foreach ([308, 426] as $code) {
  php_cli_server_start(<<<PHP
http_response_code($code);
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


HEADER
  )) {
      while (!feof($fp)) {
          echo fgets($fp);
      }
  }

  fclose($fp);
}
?>
