<?php

$xml = <<<XML
<?xml version="1.0" encoding="utf-8"?>
<params>
<param>
 <value>
  <string>Is string</string>
 </value>
</param>
</params>
XML;

$response = xmlrpc_decode($xml);
echo $response;

