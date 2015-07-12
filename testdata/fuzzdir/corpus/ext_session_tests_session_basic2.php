<?php

ob_start();

/*
 * Prototype : session.use_strict_mode=1
 * Description : Test basic functionality.
 * Source code : ext/session/session.c, ext/session/mod_files.c
 */

echo "*** Testing basic session functionality : variation2 ***\n";

$session_id = 'testid';
session_id($session_id);
$path = dirname(__FILE__);
var_dump(session_save_path($path));

echo "*** Without lazy_write ***\n";
var_dump(session_id($session_id));
var_dump(session_start(['lazy_write'=>FALSE]));
$session_id_new1 = session_id();
var_dump($session_id_new1 !== $session_id);
var_dump(session_write_close());
var_dump(session_id());

echo "*** With lazy_write ***\n";
var_dump(session_id($session_id));
var_dump(session_start(['lazy_write'=>TRUE]));
$session_id_new2 = session_id();
var_dump($session_id_new1 !== $session_id_new2);
var_dump(session_commit());
var_dump(session_id());

echo "*** Cleanup ***\n";
ini_set('session.use_strict_mode',0);
var_dump(session_id($session_id_new1));
var_dump(session_start());
var_dump(session_destroy());
var_dump(session_id($session_id_new2));
var_dump(session_start());
var_dump(session_destroy());

ob_end_flush();
?>
