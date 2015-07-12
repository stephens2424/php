<?php
$handler="flatfile";
require_once(dirname(__FILE__) .'/test.inc');
echo "database handler: $handler\n";

function check($h)
{
    if (!$h) {
        return;
    }

    foreach ($h as $key) {
        if ($key === "flatfile") {
            echo "Success: flatfile enabled\n";
        }
    }
}

echo "Test 1\n";

check(dba_handlers());

echo "Test 2\n";

check(dba_handlers(null));

echo "Test 3\n";

check(dba_handlers(1, 2));

echo "Test 4\n";

check(dba_handlers(0));

echo "Test 5 - full info\n";
$h = dba_handlers(1);
foreach ($h as $key => $val) {
    if ($key === "flatfile") {
        echo "Success: flatfile enabled\n";
    }
}

?>
