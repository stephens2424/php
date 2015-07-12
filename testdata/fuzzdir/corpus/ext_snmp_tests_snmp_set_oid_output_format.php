<?php
require_once(dirname(__FILE__).'/snmp_include.inc');

echo "Checking error handling\n";
var_dump(snmp_set_oid_output_format());
var_dump(snmp_set_oid_output_format(123));

echo "Checking working\n";
var_dump(snmp_set_oid_output_format(SNMP_OID_OUTPUT_FULL));
var_dump(snmp_set_oid_output_format(SNMP_OID_OUTPUT_NUMERIC));
?>
