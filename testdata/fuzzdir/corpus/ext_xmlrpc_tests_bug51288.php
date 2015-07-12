<?php
$method = NULL;
$req = '<?xml version="1.0"?><methodCall></methodCall>';
var_dump(xmlrpc_decode_request($req, $method));
var_dump($method);
echo "Done\n";
?>
