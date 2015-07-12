<?php
require "connect.inc";
$link = ldap_connect_and_bind($host, $port, $user, $passwd, $protocol_version);
insert_dummy_data($link, $base);
ldap_add($link, "cn=userref,$base", array(
        "objectClass" => array("extensibleObject", "referral"),
        "cn" => "userref",
        "ref" => "cn=userA,$base",
));
ldap_set_option($link, LDAP_OPT_DEREF, LDAP_DEREF_NEVER);
$result = ldap_search($link, "$base", "(cn=*)");
var_dump($ref = ldap_first_reference($link, $result));
$refs = null;
ldap_parse_reference($link, $ref, $refs);
var_dump($refs);
?>
===DONE===
