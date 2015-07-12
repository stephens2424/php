<?php

function DecodeDatetime($datetime) {
  print "\nISO 8601 datetime $datetime\n";
  $obj = xmlrpc_decode("<?xml version=\"1.0\"?><methodResponse><params><param><value><dateTime.iso8601>$datetime</dateTime.iso8601></value></param></params></methodResponse>");
  print_r($obj);
}

DecodeDatetime("20010909T01:46:40Z");
DecodeDatetime("20010909T00:46:40-01");
DecodeDatetime("2001-09-09T08:46:40+07:00");
DecodeDatetime("2001-09-08T21:46:40-0400");

?>
