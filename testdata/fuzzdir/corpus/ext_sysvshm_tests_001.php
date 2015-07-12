<?php

var_dump(ftok());
var_dump(ftok(1));
var_dump(ftok(1,1,1));

var_dump(ftok("",""));
var_dump(ftok(-1, -1));
var_dump(ftok("qwertyu","qwertyu"));

var_dump(ftok("nonexistentfile","q"));

var_dump(ftok(__FILE__,"q"));

echo "Done\n";
?>
