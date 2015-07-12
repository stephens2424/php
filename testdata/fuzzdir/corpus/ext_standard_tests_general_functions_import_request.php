<?php

var_dump(import_request_variables());
var_dump(import_request_variables(""));
var_dump(import_request_variables("", ""));

var_dump(import_request_variables("g", ""));
var_dump($a, $b, $c, $ap);

var_dump(import_request_variables("g", "g_"));
var_dump($g_a, $g_b, $g_c, $g_ap, $g_1);

var_dump(import_request_variables("GP", "i_"));
var_dump($i_a, $i_b, $i_c, $i_ap, $i_bp, $i_cp, $i_dp);

var_dump(import_request_variables("gGg", "r_"));
var_dump($r_a, $r_b, $r_c, $r_ap);

echo "Done\n";
?>
