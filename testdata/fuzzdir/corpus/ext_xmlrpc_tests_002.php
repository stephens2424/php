<?php

$r = xmlrpc_encode_request("method", array());
var_dump(xmlrpc_decode_request($r, $method));
var_dump($method);

$r = xmlrpc_encode_request("method", 1);
var_dump(xmlrpc_decode_request($r, $method));
var_dump($method);

$r = xmlrpc_encode_request("method", 'param');
var_dump(xmlrpc_decode_request($r, $method));
var_dump($method);

$r = xmlrpc_encode_request(-1, "");
var_dump(xmlrpc_decode_request($r, $method));
var_dump($method);

$r = xmlrpc_encode_request(array(), 1);
var_dump(xmlrpc_decode_request($r, $method));
var_dump($method);

echo "Done\n";
?>
