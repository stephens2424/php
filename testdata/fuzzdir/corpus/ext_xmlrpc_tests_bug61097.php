<?php
$server = xmlrpc_server_create();

$method = 'abc';
xmlrpc_server_register_introspection_callback($server, $method);
xmlrpc_server_register_method($server, 'abc', $method);

echo 'Done';
?>
