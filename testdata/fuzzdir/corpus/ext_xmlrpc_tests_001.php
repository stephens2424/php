<?php

var_dump(xmlrpc_encode_request(-1, 1));
var_dump(xmlrpc_encode_request("", 1));
var_dump(xmlrpc_encode_request(array(), 1));
var_dump(xmlrpc_encode_request(3.4, 1));

echo "Done\n";
?>
