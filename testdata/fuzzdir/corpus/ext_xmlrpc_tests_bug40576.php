<?php

var_dump(xmlrpc_encode(1.123456789));
var_dump(xmlrpc_encode(11234567891010));
var_dump(xmlrpc_encode(11234567));
var_dump(xmlrpc_encode(""));
var_dump(xmlrpc_encode("test"));
var_dump(xmlrpc_encode("1.22222222222222222222222"));

echo "Done\n";
?>
