<?php
class ancester
{
    public $ancester = 0;
	function __construct()
	{
		return $this->ancester;
	}
}
class foo extends ancester
{
    public $bar = "1";
	function __construct()
	{
		return $this->bar;
	}
}

$r = new ReflectionClass('foo');
foreach ($r->GetProperties() as $p)
{
	echo $p->getName(). " ". $p->getDeclaringClass()->getName()."\n";
}

?>
