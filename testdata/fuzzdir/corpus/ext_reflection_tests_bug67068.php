<?php
class MyClass {
    public function method() {}
}

$object = new MyClass;
$reflector = new \ReflectionMethod($object, 'method');
$closure = $reflector->getClosure($object);

$closureReflector = new \ReflectionFunction($closure);

var_dump($closureReflector->isClosure());
?>
