<?php
include "php_cli_server.inc";
php_cli_server_start();
foreach (['MKCALENDAR', 'MKCO', 'MKCOLL', 'M'] as $method) {
    $context = stream_context_create(['http' => ['method' => $method]]);
    // the following is supposed to emit a warning for unsupported methods
    file_get_contents("http://" . PHP_CLI_SERVER_ADDRESS, false, $context);
}
?>
