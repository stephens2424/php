<?php
echo wddx_serialize_value("\xfc\x63") . "\n";
echo wddx_serialize_value([ "\xfc\x63" => "foo" ]) . "\n";
?>
