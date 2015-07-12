<?php

ob_start();

/*
 * Prototype : bool session_set_save_handler(callback $open, callback $close, callback $read, callback $write, callback $destroy, callback $gc)
 * Description : Sets user-level session storage functions with validate_id() and update()
 * Source code : ext/session/session.c
 */

function noisy_gc($maxlifetime) {
	echo("GC [".$maxlifetime."]\n");
	echo gc($maxlifetime)." deleted\n";
	return true;
}

echo "*** Testing session_set_save_handler() : variation ***\n";

require_once "save_handler.inc";
$path = dirname(__FILE__);
var_dump(session_save_path($path));

echo "*** Without lazy_write ***\n";
var_dump(session_set_save_handler("open", "close", "read", "write", "destroy", "noisy_gc", "create_sid", "validate_sid", "update"));
var_dump(session_start(['lazy_write'=>FALSE]));
$session_id = session_id();
var_dump(session_id());
var_dump(session_write_close());
var_dump(session_id());

echo "*** With lazy_write ***\n";
var_dump(session_id($session_id));
var_dump(session_set_save_handler("open", "close", "read", "write", "destroy", "noisy_gc", "create_sid", "validate_sid", "update"));
var_dump(session_start(['lazy_write'=>TRUE]));
var_dump(session_commit());
var_dump(session_id());

echo "*** Cleanup ***\n";
var_dump(session_id($session_id));
var_dump(session_start());
var_dump(session_destroy());

ob_end_flush();
?>
