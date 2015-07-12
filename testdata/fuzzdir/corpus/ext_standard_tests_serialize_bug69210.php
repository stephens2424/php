<?php
class testString
{
	public $a = true;

	public function __sleep()
	{
		return array('a', '1');
	}
}

class testInteger
{
	public $a = true;

	public function __sleep()
	{
		return array('a', 1);
	}
}

$cs = new testString();
$ci = new testInteger();

$ss =  @serialize($cs);
echo $ss . "\n";

$si = @serialize($ci);
echo $si . "\n";

var_dump(unserialize($ss));
var_dump(unserialize($si));
?>
