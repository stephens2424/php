<?php
require "connect.inc";

$link = ldap_connect_and_bind($host, $port, $user, $passwd, $protocol_version);
insert_dummy_data($link, $base);
ldap_add($link, "cn=userref,$base", array(
        "objectClass" => array("extensibleObject", "referral"),
        "cn" => "userref",
        "ref" => "cn=userA,$base",
));
$result = ldap_search($link, "cn=userref,$base", "(cn=user*)");
$errcode = $dn = $errmsg = $refs =  null;
var_dump(
	ldap_parse_result($link, $result, $errcode, $dn, $errmsg, $refs),
	$errcode, $dn, $errmsg, $refs
);
?>
===DONE===
