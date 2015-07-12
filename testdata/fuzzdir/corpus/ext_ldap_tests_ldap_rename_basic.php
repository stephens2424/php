<?php
require "connect.inc";

$link = ldap_connect_and_bind($host, $port, $user, $passwd, $protocol_version);
insert_dummy_data($link, $base);
var_dump(
	ldap_rename($link, "cn=userA,$base", "cn=userZ", "$base", true)
);
$result = ldap_search($link, "$base", "(cn=userA)", array("cn", "sn"));
$result = ldap_search($link, "$base", "(cn=userZ)", array("cn", "sn"));
var_dump(ldap_get_entries($link, $result));
?>
===DONE===
