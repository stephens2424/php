<?php
require_once(dirname(__FILE__).'/snmp_include.inc');

echo "Checking error handling\n";
var_dump(snmp_get_quick_print('noarg'));
var_dump(snmp_set_quick_print('noarg'));
var_dump(snmp_set_quick_print());

echo "Checking working\n";
var_dump(snmp_get_quick_print());
snmp_set_quick_print(false);
var_dump(snmp_get_quick_print());
snmp_set_quick_print(true);
var_dump(snmp_get_quick_print());

?>
