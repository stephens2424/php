<?php
var_dump(mb_ereg_replace("C?$", "Z", "ABC"));
var_dump(preg_replace("/C?$/", "Z", "ABC"));
var_dump(mb_ereg_replace("C*$", "Z", "ABC"));
var_dump(preg_replace("/C*$/", "Z", "ABC"));
?>
