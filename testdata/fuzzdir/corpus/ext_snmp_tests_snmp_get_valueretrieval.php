<?php
require_once(dirname(__FILE__).'/snmp_include.inc');

echo "Checking error handling\n";
var_dump(snmp_get_valueretrieval('noarg'));
var_dump(snmp_set_valueretrieval());
var_dump(snmp_set_valueretrieval('noarg'));
var_dump(snmp_set_valueretrieval(67));

echo "Checking working\n";
var_dump(snmp_get_valueretrieval());
snmp_set_valueretrieval(SNMP_VALUE_LIBRARY);
var_dump(snmp_get_valueretrieval() === SNMP_VALUE_LIBRARY);
snmp_set_valueretrieval(SNMP_VALUE_PLAIN);
var_dump(snmp_get_valueretrieval() === SNMP_VALUE_PLAIN);
snmp_set_valueretrieval(SNMP_VALUE_OBJECT);
var_dump(snmp_get_valueretrieval() === SNMP_VALUE_OBJECT);
snmp_set_valueretrieval(SNMP_VALUE_PLAIN|SNMP_VALUE_OBJECT);
var_dump(snmp_get_valueretrieval() === (SNMP_VALUE_PLAIN|SNMP_VALUE_OBJECT));
snmp_set_valueretrieval(SNMP_VALUE_LIBRARY|SNMP_VALUE_OBJECT);
var_dump(snmp_get_valueretrieval() === (SNMP_VALUE_LIBRARY|SNMP_VALUE_OBJECT));

?>
