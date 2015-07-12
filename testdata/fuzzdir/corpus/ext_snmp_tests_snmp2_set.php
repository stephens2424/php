<?php
require_once(dirname(__FILE__).'/snmp_include.inc');

//EXPECTF format is quickprint OFF
snmp_set_quick_print(false);
snmp_set_valueretrieval(SNMP_VALUE_PLAIN);

echo "Check error handing\n";
echo "4args (5 needed)\n";
$z = snmp2_set($hostname, $communityWrite, 'SNMPv2-MIB::sysLocation.0');
var_dump($z);

echo "No type & no value (timeout & retries instead)\n";
$z = snmp2_set($hostname, $communityWrite, 'SNMPv2-MIB::sysLocation.0', $timeout, $retries);
var_dump($z);

echo "No value (timeout instead), retries instead of timeout\n";
$z = snmp2_set($hostname, $communityWrite, 'SNMPv2-MIB::sysLocation.0', 'q', $timeout, $retries);
var_dump($z);

echo "Bogus OID\n";
$z = snmp2_set($hostname, $communityWrite, '.1.3.6.777.888.999.444.0', 's', 'bbb', $timeout, $retries);
var_dump($z);

echo "Checking working\n";
$oid1 = 'SNMPv2-MIB::sysContact.0';
$oldvalue1 = snmpget($hostname, $communityWrite, $oid1, $timeout, $retries);
$newvalue1 = $oldvalue1 . '0';
$oid2 = 'SNMPv2-MIB::sysLocation.0';
$oldvalue2 = snmpget($hostname, $communityWrite, $oid1, $timeout, $retries);
$newvalue2 = $oldvalue2 . '0';

echo "Single OID\n";
$z = snmp2_set($hostname, $communityWrite, $oid1, 's', $newvalue1, $timeout, $retries);
var_dump($z);
var_dump((snmpget($hostname, $communityWrite, $oid1, $timeout, $retries) === $newvalue1));
$z = snmp2_set($hostname, $communityWrite, $oid1, 's', $oldvalue1, $timeout, $retries);
var_dump($z);
var_dump((snmpget($hostname, $communityWrite, $oid1, $timeout, $retries) === $oldvalue1));

echo "Multiple OID\n";
$z = snmp2_set($hostname, $communityWrite, array($oid1, $oid2), array('s','s'), array($newvalue1, $newvalue2), $timeout, $retries);
var_dump($z);
var_dump((snmpget($hostname, $communityWrite, $oid1, $timeout, $retries) === $newvalue1));
var_dump((snmpget($hostname, $communityWrite, $oid2, $timeout, $retries) === $newvalue2));
$z = snmp2_set($hostname, $communityWrite, array($oid1, $oid2), array('s','s'), array($oldvalue1, $oldvalue2), $timeout, $retries);
var_dump($z);
var_dump((snmpget($hostname, $communityWrite, $oid1, $timeout, $retries) === $oldvalue1));
var_dump((snmpget($hostname, $communityWrite, $oid2, $timeout, $retries) === $oldvalue2));

echo "Multiple OID, single type & value\n";
$z = snmp2_set($hostname, $communityWrite, array($oid1, $oid2), 's', $newvalue1, $timeout, $retries);
var_dump($z);
var_dump((snmpget($hostname, $communityWrite, $oid1, $timeout, $retries) === $newvalue1));
var_dump((snmpget($hostname, $communityWrite, $oid2, $timeout, $retries) === $newvalue1));
$z = snmp2_set($hostname, $communityWrite, array($oid1, $oid2), array('s','s'), array($oldvalue1, $oldvalue2), $timeout, $retries);
var_dump($z);
var_dump((snmpget($hostname, $communityWrite, $oid1, $timeout, $retries) === $oldvalue1));
var_dump((snmpget($hostname, $communityWrite, $oid2, $timeout, $retries) === $oldvalue2));

echo "Multiple OID, single type, multiple value\n";
$z = snmp2_set($hostname, $communityWrite, array($oid1, $oid2), 's', array($newvalue1, $newvalue2), $timeout, $retries);
var_dump($z);
var_dump((snmpget($hostname, $communityWrite, $oid1, $timeout, $retries) === $newvalue1));
var_dump((snmpget($hostname, $communityWrite, $oid2, $timeout, $retries) === $newvalue2));
$z = snmp2_set($hostname, $communityWrite, array($oid1, $oid2), array('s','s'), array($oldvalue1, $oldvalue2), $timeout, $retries);
var_dump($z);
var_dump((snmpget($hostname, $communityWrite, $oid1, $timeout, $retries) === $oldvalue1));
var_dump((snmpget($hostname, $communityWrite, $oid2, $timeout, $retries) === $oldvalue2));


echo "More error handing\n";
echo "Single OID, single type in array, single value\n";
$z = snmp2_set($hostname, $communityWrite, $oid1, array('s'), $newvalue1, $timeout, $retries);
var_dump($z);
var_dump((snmpget($hostname, $communityWrite, $oid1, $timeout, $retries) === $oldvalue1));
var_dump((snmpget($hostname, $communityWrite, $oid2, $timeout, $retries) === $oldvalue2));

echo "Single OID, single type, single value in array\n";
$z = snmp2_set($hostname, $communityWrite, $oid1, 's', array($newvalue1), $timeout, $retries);
var_dump($z);
var_dump((snmpget($hostname, $communityWrite, $oid1, $timeout, $retries) === $oldvalue1));
var_dump((snmpget($hostname, $communityWrite, $oid2, $timeout, $retries) === $oldvalue2));

echo "Multiple OID, 1st wrong type\n";
$z = snmp2_set($hostname, $communityWrite, array($oid1, $oid2), array('sw','s'), array($newvalue1, $newvalue2), $timeout, $retries);
var_dump($z);
var_dump((snmpget($hostname, $communityWrite, $oid1, $timeout, $retries) === $oldvalue1));
var_dump((snmpget($hostname, $communityWrite, $oid2, $timeout, $retries) === $oldvalue2));

echo "Multiple OID, 2nd wrong type\n";
$z = snmp2_set($hostname, $communityWrite, array($oid1, $oid2), array('s','sb'), array($newvalue1, $newvalue2), $timeout, $retries);
var_dump($z);
var_dump((snmpget($hostname, $communityWrite, $oid1, $timeout, $retries) === $oldvalue1));
var_dump((snmpget($hostname, $communityWrite, $oid2, $timeout, $retries) === $oldvalue2));

echo "Multiple OID, single type in array, multiple value\n";
$z = snmp2_set($hostname, $communityWrite, array($oid1, $oid2), array('s'), array($newvalue1, $newvalue2), $timeout, $retries);
var_dump($z);
var_dump((snmpget($hostname, $communityWrite, $oid1, $timeout, $retries) === $oldvalue1));
var_dump((snmpget($hostname, $communityWrite, $oid2, $timeout, $retries) === $oldvalue2));

echo "Multiple OID & type, singe value in array\n";
$z = snmp2_set($hostname, $communityWrite, array($oid1, $oid2), array('s', 's'), array($newvalue1), $timeout, $retries);
var_dump($z);
var_dump((snmpget($hostname, $communityWrite, $oid1, $timeout, $retries) === $oldvalue1));
var_dump((snmpget($hostname, $communityWrite, $oid2, $timeout, $retries) === $oldvalue2));

echo "Multiple OID, 1st bogus, single type, multiple value\n";
$z = snmp2_set($hostname, $communityWrite, array($oid1 . '44.55.66.77', $oid2), 's', array($newvalue1, $newvalue2), $timeout, $retries);
var_dump($z);
var_dump((snmpget($hostname, $communityWrite, $oid1, $timeout, $retries) === $oldvalue1));
var_dump((snmpget($hostname, $communityWrite, $oid2, $timeout, $retries) === $oldvalue2));

echo "Multiple OID, 2nd bogus, single type, multiple value\n";
$z = snmp2_set($hostname, $communityWrite, array($oid1, $oid2 . '44.55.66.77'), 's', array($newvalue1, $newvalue2), $timeout, $retries);
var_dump($z);
var_dump((snmpget($hostname, $communityWrite, $oid1, $timeout, $retries) === $oldvalue1));
var_dump((snmpget($hostname, $communityWrite, $oid2, $timeout, $retries) === $oldvalue2));

echo "Multiple OID, single multiple type (1st bogus), multiple value\n";
$z = snmp2_set($hostname, $communityWrite, array($oid1, $oid2), array('q', 's'), array($newvalue1, $newvalue2), $timeout, $retries);
var_dump($z);
var_dump((snmpget($hostname, $communityWrite, $oid1, $timeout, $retries) === $oldvalue1));
var_dump((snmpget($hostname, $communityWrite, $oid2, $timeout, $retries) === $oldvalue2));

echo "Multiple OID, single multiple type (2nd bogus), multiple value\n";
$z = snmp2_set($hostname, $communityWrite, array($oid1, $oid2), array('s', 'w'), array($newvalue1, $newvalue2), $timeout, $retries);
var_dump($z);
var_dump((snmpget($hostname, $communityWrite, $oid1, $timeout, $retries) === $oldvalue1));
var_dump((snmpget($hostname, $communityWrite, $oid2, $timeout, $retries) === $oldvalue2));

?>
