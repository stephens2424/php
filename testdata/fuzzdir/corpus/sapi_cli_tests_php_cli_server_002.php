<?php
include "php_cli_server.inc";
php_cli_server_start('var_dump($_SERVER["DOCUMENT_ROOT"], $_SERVER["SERVER_SOFTWARE"], $_SERVER["SERVER_NAME"], $_SERVER["SERVER_PORT"]);');
var_dump(file_get_contents("http://" . PHP_CLI_SERVER_ADDRESS));
?>
