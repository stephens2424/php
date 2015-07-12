<?php
require_once(dirname(__FILE__).'/snmp_include.inc');

//EXPECTF format is quickprint OFF
snmp_set_quick_print(false);
snmp_set_valueretrieval(SNMP_VALUE_PLAIN);

echo "Checking error handling\n";
var_dump(snmprealwalk($hostname, $community, '.1.3.6.1.2.1.1', ''));
var_dump(snmprealwalk($hostname, $community, '.1.3.6.1.2.1.1', $timeout, ''));

echo "Checking working\n";
echo "Single OID\n";
$return = snmprealwalk($hostname, $community, '.1.3.6.1.2.1.1', $timeout, $retries);
var_dump(gettype($return));
var_dump(sizeof($return));
var_dump(key($return));
var_dump(array_shift($return));

echo "Single OID in array\n";
$return = snmprealwalk($hostname, $community, array('.1.3.6.1.2.1.1'), $timeout, $retries);
var_dump(gettype($return));
var_dump(sizeof($return));
var_dump(key($return));
var_dump(array_shift($return));

echo "More error handling\n";
echo "Multiple correct OID\n";
$return = snmprealwalk($hostname, $community, array('.1.3.6.1.2.1.1', '.1.3.6'), $timeout, $retries);
var_dump($return);

echo "Multiple OID with wrong OID\n";
$return = snmprealwalk($hostname, $community, array('.1.3.6.1.2.1.1', '.1.3.6...1'), $timeout, $retries);
var_dump($return);
$return = snmprealwalk($hostname, $community, array('.1.3.6...1', '.1.3.6.1.2.1.1'), $timeout, $retries);
var_dump($return);

echo "Single nonexisting OID\n";
$return = snmprealwalk($hostname, $community, array('.1.3.6.99999.0.99999.111'), $timeout, $retries);
var_dump($return);

?>
