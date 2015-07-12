<?php
$handler="db4";
require_once(dirname(__FILE__) .'/test.inc');
echo "database handler: $handler\n";

function check($h)
{
    if (!$h) {
        return;
    }

    foreach ($h as $key) {
        if ($key === "db4") {
            echo "Success: db4 enabled\n";
        }
    }
}

echo "Test 1\n";

check(dba_handlers());

echo "Test 2 - full info\n";
$h = dba_handlers(1);
foreach ($h as $key => $val) {
    if ($key === "db4") {
        echo "$val\n";
    }
}

?>
