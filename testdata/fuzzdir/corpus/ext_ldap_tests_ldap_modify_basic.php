<?php
require "connect.inc";

$link = ldap_connect_and_bind($host, $port, $user, $passwd, $protocol_version);
insert_dummy_data($link, $base);

$entry = array(
	"objectClass"	=> array(
		"top",
		"organization"),
	"o"		=> "test",
	"description"	=> "Domain description",
);

var_dump(
	ldap_modify($link, "o=test,$base", $entry),
	ldap_get_entries(
		$link,
		ldap_search($link, "$base", "(Description=Domain description)")
	)
);
?>
===DONE===
