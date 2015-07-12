<?php
require_once(dirname(__FILE__).'/snmp_include.inc');

//EXPECTF format is quickprint OFF
snmp_set_quick_print(false);
snmp_set_valueretrieval(SNMP_VALUE_PLAIN);

$session = new SNMP(SNMP::VERSION_3, $hostname, $user_noauth, $timeout, $retries);
$session->setSecurity('noAuthNoPriv');

#echo "Checking error handling\n";
var_dump($session->setSecurity());
var_dump($session->setSecurity(''));
var_dump($session->setSecurity('bugusPriv'));
var_dump($session->setSecurity('authNoPriv', 'TTT'));
var_dump($session->setSecurity('authNoPriv', 'MD5', ''));
var_dump($session->setSecurity('authNoPriv', 'MD5', 'te'));
var_dump($session->setSecurity('authPriv', 'MD5', $auth_pass, 'BBB'));
var_dump($session->setSecurity('authPriv', 'MD5', $auth_pass, 'AES', ''));
var_dump($session->setSecurity('authPriv', 'MD5', $auth_pass, 'AES', 'ty'));
var_dump($session->setSecurity('authPriv', 'MD5', $auth_pass, 'AES', 'test12345', 'context', 'dsa'));

var_dump($session->close());

?>
