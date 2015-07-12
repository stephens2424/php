<?php
include "php_cli_server.inc";
php_cli_server_start('print_r($_REQUEST); $_REQUEST["foo"] = "bar"; return FALSE;');
$doc_root = __DIR__;
file_put_contents($doc_root . '/request.php', '<?php print_r($_REQUEST); ?>');

list($host, $port) = explode(':', PHP_CLI_SERVER_ADDRESS);
$port = intval($port)?:80;

$fp = fsockopen($host, $port, $errno, $errstr, 0.5);
if (!$fp) {
  die("connect failed");
}

if(fwrite($fp, <<<HEADER
POST /request.php HTTP/1.1
Host: {$host}
Content-Type: application/x-www-form-urlencoded
Content-Length: 3

a=b
HEADER
)) {
	while (!feof($fp)) {
		echo fgets($fp);
	}
}

fclose($fp);
@unlink($doc_root . '/request.php');

?>
