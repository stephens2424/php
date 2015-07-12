<?php
require_once(dirname(__FILE__).'/snmp_include.inc');

echo "Checking error handling\n";
var_dump(snmp_read_mib());
var_dump(snmp_read_mib(dirname(__FILE__).'/cannotfindthisfile'));

echo "Checking working\n";
var_dump(snmp_read_mib($mibdir . '/SNMPv2-MIB.txt'));

?>
