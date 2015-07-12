<?php

$d = date(DATE_ISO8601);
xmlrpc_set_type($d, 'datetime');
echo xmlrpc_encode_request('method.call', array('date' => $d));

$d = '2008-01-01 20:00:00';
xmlrpc_set_type($d, 'datetime');
echo xmlrpc_encode_request('method.call', array('date' => $d));

?>
