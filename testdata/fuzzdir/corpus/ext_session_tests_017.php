<?php

error_reporting(E_ALL);

class Kill {
	function __construct() {
		global $HTTP_SESSION_VARS;
		session_start();
	}
}
$k = new Kill();

print "I live\n";
?>
