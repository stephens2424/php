<?php

var_dump(tidy_get_opt_doc(new tidy, 'some_bogus_cfg'));

$t = new tidy;
var_dump($t->getOptDoc('ncr'));
var_dump(strlen(tidy_get_opt_doc($t, 'wrap')) > 99);
?>
