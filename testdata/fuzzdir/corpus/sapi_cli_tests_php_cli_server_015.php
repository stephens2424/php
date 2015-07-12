<?php
include "php_cli_server.inc";
php_cli_server_start('require("syntax_error.php");');
$dir = realpath(dirname(__FILE__));

file_put_contents($dir . "/syntax_error.php", "<?php non_exists_function(); ?>");

list($host, $port) = explode(':', PHP_CLI_SERVER_ADDRESS);
$port = intval($port)?:80;
$output = '';

$fp = fsockopen($host, $port, $errno, $errstr, 0.5);
if (!$fp) {
  die("connect failed");
}

if(fwrite($fp, <<<HEADER
GET /index.php HTTP/1.1
Host: {$host}


HEADER
)) {
	while (!feof($fp)) {
		$output .= fgets($fp);
	}
}
echo $output;
@unlink($dir . "/syntax_error.php");
fclose($fp);
?>
