<?php

$xml = <<<XML
<?xml version="1.0" encoding="UTF-8"?>
<aa>
    <bs>
        <b>b</b>
    </bs>
    <cs><c>b</c></cs>
    <ds><d id="d"></d></ds>
    <es>
        <e id="e"></e>
    </es>
    <fs><f id="f"></f><f id="f"></f></fs>
</aa>
XML;
$sxe = simplexml_load_string($xml);
print_r($sxe);

?>
