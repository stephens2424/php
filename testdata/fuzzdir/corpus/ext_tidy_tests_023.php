<?php

//test leaks here:
new tidy();
var_dump(new tidy());

echo "-------\n";

$tidy = new tidy();
$tidy->parseString('<html><?php echo "xpto;" ?></html>');

var_dump(tidy_get_root($tidy)->child[0]->isHtml());
var_dump(tidy_get_root($tidy)->child[0]->child[0]->isPHP());
var_dump(tidy_get_root($tidy)->child[0]->child[0]->isAsp());
var_dump(tidy_get_root($tidy)->child[0]->child[0]->isJste());
var_dump(tidy_get_root($tidy)->child[0]->child[0]->type === TIDY_NODETYPE_PHP);

var_dump(tidy_get_root($tidy)->child[0]->hasChildren());
var_dump(tidy_get_root($tidy)->child[0]->child[0]->hasChildren());

?>
