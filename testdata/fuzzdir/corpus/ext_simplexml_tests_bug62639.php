<?php

class A extends SimpleXMLElement
{
}

$xml1 = <<<XML
<?xml version="1.0"?>
<a>
    <b>
        <c>
            <value attr="Some Attr">Some Value</value>
        </c>
    </b>
</a>
XML;

$a1 = new A($xml1);

foreach ($a1->b->c->children() as $key => $value) {
    var_dump($value);
}

$xml2 = <<<XML
<?xml version="1.0"?>
<a>
    <b>
        <c><value attr="Some Attr">Some Value</value></c>
    </b>
</a>
XML;

$a2 = new A($xml2);

foreach ($a2->b->c->children() as $key => $value) {
    var_dump($value);
}?>
