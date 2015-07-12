<?php

$options = array ();
$request = xmlrpc_encode_request ("system.describeMethods", $options);
$server = xmlrpc_server_create ();

xmlrpc_server_register_introspection_callback($server, 1);
xmlrpc_server_register_introspection_callback($server, array('foo', 'bar'));

$options = array ('output_type' => 'xml', 'version' => 'xmlrpc');
xmlrpc_server_call_method ($server, $request, NULL, $options);

?>
