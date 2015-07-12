<?php
header("Vary: Cookie");
ob_start("ob_gzhandler");

// run-tests cannot test for a multiple Vary header
ob_flush();
print_r(headers_list());

?>
Done
