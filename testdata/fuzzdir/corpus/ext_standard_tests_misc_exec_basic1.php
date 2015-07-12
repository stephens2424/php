<?php
$cmd = "echo abc\n\0command";
var_dump(exec($cmd, $output));
var_dump($output);
var_dump(system($cmd));
var_dump(passthru($cmd));
?>
