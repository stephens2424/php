<?php
require_once(dirname(__FILE__).'/snmp_include.inc');

//EXPECTF format is quickprint OFF
snmp_set_quick_print(false);
snmp_set_valueretrieval(SNMP_VALUE_PLAIN);

echo "Checking error handling\n";
var_dump(snmp2_get($hostname, $community, '.1.3.6.1.2.1.1.1.0', ''));
var_dump(snmp2_get($hostname, $community, '.1.3.6.1.2.1.1.1.0', $timeout, ''));
echo "Empty OID array\n";
var_dump(snmp2_get($hostname, $community, array(), $timeout, $retries));

echo "Checking working\n";
echo "Single OID\n";
var_dump(snmp2_get($hostname, $community, '.1.3.6.1.2.1.1.1.0', $timeout, $retries));
echo "Single OID in array\n";
var_dump(snmp2_get($hostname, $community, array('.1.3.6.1.2.1.1.1.0'), $timeout, $retries));
echo "Multiple OID\n";
var_dump(snmp2_get($hostname, $community, array('.1.3.6.1.2.1.1.1.0', '.1.3.6.1.2.1.1.3.0'), $timeout, $retries));

echo "More error handling\n";
echo "Single OID\n";
var_dump(snmp2_get($hostname, $community, '.1.3.6.1.2..1.1.1.0', $timeout, $retries));
echo "Single OID in array\n";
var_dump(snmp2_get($hostname, $community, array('.1.3.6.1.2.1...1.1.0'), $timeout, $retries));
echo "Multiple OID\n";
var_dump(snmp2_get($hostname, $community, array('.1.3.6.1.2.1...1.1.0', '.1.3.6.1.2.1.1.3.0'), $timeout, $retries));

echo "noSuchName checks\n";
echo "Single OID\n";
var_dump(snmp2_get($hostname, $community, '.1.3.6.1.2.1.1.1.110', $timeout, $retries));
echo "Single OID in array\n";
var_dump(snmp2_get($hostname, $community, array('.1.3.6.1.2.1.1.1.110'), $timeout, $retries));
echo "Multiple OID\n";
var_dump(snmp2_get($hostname, $community, array('.1.3.6.1.2.1.1.1.0', '.1.3.6.1.2.1.1.3.220'), $timeout, $retries));


?>
