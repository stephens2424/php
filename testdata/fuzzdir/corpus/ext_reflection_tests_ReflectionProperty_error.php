<?php

class C {
    public static $p;
}

try {
	new ReflectionProperty();
} catch (TypeError $re) {
	echo "Ok - ".$re->getMessage().PHP_EOL;
}
try {
	new ReflectionProperty('C::p');
} catch (TypeError $re) {
	echo "Ok - ".$re->getMessage().PHP_EOL;
}

try {
	new ReflectionProperty('C', 'p', 'x');
} catch (TypeError $re) {
	echo "Ok - ".$re->getMessage().PHP_EOL;
}


$rp = new ReflectionProperty('C', 'p');
var_dump($rp->getName(1));
var_dump($rp->isPrivate(1));
var_dump($rp->isProtected(1));
var_dump($rp->isPublic(1));
var_dump($rp->isStatic(1));
var_dump($rp->getModifiers(1));
var_dump($rp->isDefault(1));

?>
