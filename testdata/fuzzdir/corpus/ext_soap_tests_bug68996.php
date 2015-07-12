<?php
$s = new SoapServer(NULL, [
    'uri' => 'http://foo',
]);

function foo() {
  return new SoapFault("\xfc\x63", "some msg");
}
$s->addFunction("foo");

// soap 1.1
$HTTP_RAW_POST_DATA = <<<EOF
<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
  <SOAP-ENV:Body>
    <SOAP-ENV:foo />
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>
EOF;
$s->handle($HTTP_RAW_POST_DATA);

// soap 1.2
$HTTP_RAW_POST_DATA = <<<EOF
<?xml version="1.0" encoding="UTF-8"?>
<env:Envelope xmlns:env="http://www.w3.org/2003/05/soap-envelope">
  <env:Body>
    <env:foo />
  </env:Body>
</env:Envelope>
EOF;
$s->handle($HTTP_RAW_POST_DATA);
?>
