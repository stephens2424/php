<?php
class SOAPStruct {
    function __construct($s, $i, $f) {
        $this->varString = $s;
        $this->varInt = $i;
        $this->varFloat = $f;
    }
}
class BaseStruct {
    function __construct($f, $s) {
        $this->structMessage = $f;
        $this->shortMessage = $s;
    }
}
class ExtendedStruct extends BaseStruct {
    function __construct($f, $s, $x1, $x2, $x3) {
        parent::__construct($f,$s);
        $this->stringMessage = $x1;
        $this->intMessage = $x2;
        $this->anotherIntMessage = $x3;
    }
}
$struct = new ExtendedStruct(new SOAPStruct("a1",11,12.345),12,"arg",-3,5);
$client = new SoapClient(dirname(__FILE__)."/round4_groupH_complex_doclit.wsdl",array("trace"=>1,"exceptions"=>0));
$client->echoExtendedStructFault($struct);
echo $client->__getlastrequest();
$HTTP_RAW_POST_DATA = $client->__getlastrequest();
include("round4_groupH_complex_doclit.inc");
echo "ok\n";
?>
