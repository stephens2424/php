<?php
header("A: first");
header("A: second", TRUE);
$headers1 = headers_list();
header("A: third", FALSE);
$headers2 = headers_list();
header_remove("A");
$headers3 = headers_list();
print_r($headers1);
print_r($headers2);
print_r($headers3);
