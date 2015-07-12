<?php
$p12 = file_get_contents(__DIR__.'/p12_with_extra_certs.p12');

$result = openssl_pkcs12_read($p12, $cert_data, 'qwerty');
var_dump($result);
var_dump(openssl_error_string());
?>
