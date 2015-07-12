<?php
require_once(dirname(__FILE__).'/snmp_include.inc');

# hostname variable was modified inline in netsnmp_session_init()
# Should be checked with IPv6 since IPv4 processing code do not alter pointer position

//EXPECTF format is quickprint OFF
snmp_set_quick_print(false);
snmp_set_valueretrieval(SNMP_VALUE_PLAIN);

$checkvar = "$hostname6_port";

var_dump(snmpget($checkvar, $community, '.1.3.6.1.2.1.1.1.0'));
var_dump(($checkvar === $hostname6_port));
var_dump(snmpget($checkvar, $community, '.1.3.6.1.2.1.1.1.0'));
var_dump(($checkvar === $hostname6_port));
var_dump(snmpget($checkvar, $community, '.1.3.6.1.2.1.1.1.0'));
var_dump(($checkvar === $hostname6_port));
?>
