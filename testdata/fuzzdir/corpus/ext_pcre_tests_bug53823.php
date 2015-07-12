<?php
var_dump(preg_replace('/[^\pL\pM]*/iu', '', 'áéíóú'));
// invalid UTF-8
var_dump(preg_replace('/[^\pL\pM]*/iu', '', "\xFCáéíóú"));
var_dump(preg_replace('/[^\pL\pM]*/iu', '', "áéíóú\xFC"));
?>
