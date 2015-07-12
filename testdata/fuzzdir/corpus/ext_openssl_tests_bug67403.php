<?php
$r = openssl_x509_parse(file_get_contents(__DIR__.'/bug64802.pem'));
var_dump($r['signatureTypeSN']);
var_dump($r['signatureTypeLN']);
var_dump($r['signatureTypeNID']);

$r = openssl_x509_parse(file_get_contents(__DIR__.'/bug37820cert.pem'));
var_dump($r['signatureTypeSN']);
var_dump($r['signatureTypeLN']);
var_dump($r['signatureTypeNID']);
