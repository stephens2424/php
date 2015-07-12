<?php
$url = "file:///etc/passwd\0http://google.com";
$ch = curl_init();
var_dump(curl_setopt($ch, CURLOPT_URL, $url));
?>
Done
