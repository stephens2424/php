<?php
$mh = curl_multi_init();
$array = array($mh);
$array[] = &$array;

curl_multi_add_handle($mh, curl_init());
curl_multi_add_handle($mh, curl_init());
curl_multi_add_handle($mh, curl_init());
curl_multi_add_handle($mh, curl_init());
echo "okey";
?>
