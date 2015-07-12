<?php

$response = xmlrpc_encode(3.24234);
echo $response;

$response = xmlrpc_encode(-3.24234);
echo $response;

$response = xmlrpc_encode('Is string');
echo $response;

