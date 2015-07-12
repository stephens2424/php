<?php
include "php_cli_server.inc";
php_cli_server_start('chdir(__DIR__); echo "okey";');
var_dump(file_get_contents("http://" . PHP_CLI_SERVER_ADDRESS));
var_dump(file_get_contents("http://" . PHP_CLI_SERVER_ADDRESS));
?>
