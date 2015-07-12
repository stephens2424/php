<?php
mkdir("bug45181_x");
var_dump(is_dir("bug45181_x"));
chdir("bug45181_x");
var_dump(is_dir("bug45181_x"));
?>
