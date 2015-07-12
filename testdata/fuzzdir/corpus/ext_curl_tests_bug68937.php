<?php

$ch = curl_init('http://www.google.com/');
curl_setopt_array($ch, array(
	CURLOPT_HEADER => false,
	CURLOPT_RETURNTRANSFER => true,
	CURLOPT_POST => true,
	CURLOPT_INFILESIZE => 1,
	CURLOPT_HTTPHEADER => array(
		'Content-Length: 1',
	),
	CURLOPT_READFUNCTION => 'curl_read'
));

function curl_read($ch, $fp, $len) {
	var_dump($fp);
	exit;
}

curl_exec($ch);
curl_close($ch);
?>
