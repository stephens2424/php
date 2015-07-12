<?php

$options = array ();
$request = xmlrpc_encode_request ("system.describeMethods", $options);
$server = xmlrpc_server_create ();


function foo() { return 11111; }

class bar {
	static public function test() {
		return 'foo';
	}
}

xmlrpc_server_register_introspection_callback($server, 'foobar');
xmlrpc_server_register_introspection_callback($server, array('bar', 'test'));
xmlrpc_server_register_introspection_callback($server, array('foo', 'bar'));

$options = array ('output_type' => 'xml', 'version' => 'xmlrpc');
xmlrpc_server_call_method ($server, $request, NULL, $options);

?>
