<?php
$response='<?xml version="1.0"?>
<methodResponse>
  <params>
    <param>
      <value>
        <struct>
          <member>
            <name>50</name>
            <value><string>0.29</string></value>
          </member>
        </struct>
      </value>
    </param>
  </params>
</methodResponse>';

$retval=xmlrpc_decode($response);
var_dump($retval);
var_dump($retval["50"]);
var_dump($retval[50]);

$response='<?xml version="1.0"?>
<methodResponse>
  <params>
    <param>
      <value>
        <struct>
          <member>
            <name>0</name>
            <value><string>0.29</string></value>
          </member>
        </struct>
      </value>
    </param>
  </params>
</methodResponse>';

$retval=xmlrpc_decode($response);
var_dump($retval);
var_dump($retval["0"]);
var_dump($retval[0]);

echo "Done\n";
?>
