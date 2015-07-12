<?php
include "php_cli_server.inc";
php_cli_server_start();
var_dump(file_get_contents("http://" . PHP_CLI_SERVER_ADDRESS));
?>
