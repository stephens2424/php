<?php
class MyTidy extends tidy
{
            // foo
}

function doSomething(MyTidy $o)
{
            var_dump($o);
}

$o = new MyTidy();
var_dump($o instanceof MyTidy);
doSomething($o);
?>
