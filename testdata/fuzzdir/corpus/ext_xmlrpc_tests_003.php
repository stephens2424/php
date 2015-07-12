<?php

$params = array(
	"one" => "red",
	"two" => "blue",
	"three" => "green"
);

$response = xmlrpc_encode($params);
echo $response;

$params = array(
	"red",
	"blue",
	"green"
);

$response = xmlrpc_encode($params);
echo $response;

$params = array(
	0 => "red",
	1 => "blue",
	3 => "green"
);

$response = xmlrpc_encode($params);
echo $response;

