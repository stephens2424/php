<?php

ob_start();

/*
 * Prototype : session.use_strict_mode=0
 * Description : Test basic functionality.
 * Source code : ext/session/session.c, ext/session/mod_files.c
 */

echo "*** Testing basic session functionality : variation1 ***\n";

$session_id = 'testid';
session_id($session_id);
$path = dirname(__FILE__);
var_dump(session_save_path($path));

echo "*** Without lazy_write ***\n";
var_dump(session_id($session_id));
$config = ['lazy_write'=>FALSE];
var_dump(session_start($config));
var_dump($config);
var_dump(session_write_close());
var_dump(session_id());

echo "*** With lazy_write ***\n";
var_dump(session_id($session_id));
var_dump(session_start(['lazy_write'=>TRUE]));
var_dump(session_commit());
var_dump(session_id());

echo "*** Cleanup ***\n";
var_dump(session_id($session_id));
var_dump(session_start());
var_dump(session_destroy());

ob_end_flush();
?>
