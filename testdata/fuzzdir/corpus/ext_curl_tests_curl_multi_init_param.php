<?php
/* Prototype         : resource curl_multi_init(void)
 * Description       : Returns a new cURL multi handle
 * Source code       : ext/curl/multi.c
 * Test documentation:  http://wiki.php.net/qa/temp/ext/curl
 */

// start testing

//create the multiple cURL handle
$mh = curl_multi_init('test');
var_dump($mh);

?>
===DONE===
