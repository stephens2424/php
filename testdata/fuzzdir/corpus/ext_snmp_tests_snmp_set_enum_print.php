<?php
require_once(dirname(__FILE__).'/snmp_include.inc');

echo "Checking error handling\n";
var_dump(snmp_set_enum_print());

echo "Checking working\n";
var_dump(snmp_set_enum_print(0));
var_dump(snmp_set_enum_print(1));
?>
