<?php

$d = '6-01-01 20:00:00';
xmlrpc_set_type($d, 'datetime');
var_dump($d);
$datetime = "2001-0-08T21:46:40-0400";
$obj = xmlrpc_decode("<?xml version=\"1.0\"?><methodResponse><params><param><value><dateTime.iso8601>$datetime</dateTime.iso8601></value></param></params></methodResponse>");
print_r($obj);

$datetime = "34770-0-08T21:46:40-0400";
$obj = xmlrpc_decode("<?xml version=\"1.0\"?><methodResponse><params><param><value><dateTime.iso8601>$datetime</dateTime.iso8601></value></param></params></methodResponse>");
print_r($obj);

echo "Done\n";
?>
