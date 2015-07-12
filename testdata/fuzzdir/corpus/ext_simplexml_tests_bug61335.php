<?php
$rec1 = simplexml_load_string("<foo><bar>aa</bar>\n</foo>");
$rec2 = simplexml_load_string("<foo><bar>aa</bar></foo>");

if ($rec1->bar[0])      echo "NONEMPTY1\n";
if ($rec1->bar[0] . "") echo "NONEMPTY2\n";
if ($rec2->bar[0])      echo "NONEMPTY3\n";
?>
