<?php
class Foo {

    private static $instance;
    public static function getInstance() {
        if (isset(self::$instance)) {
            return self::$instance;
        }
        return self::$instance = new self();
    }

    public function __destruct() {
        Bar::getInstance();
    }
}

class Bar {

    private static $instance;
    public static function getInstance() {
        if (isset(self::$instance)) {
            return self::$instance;
        }
        return self::$instance = new self();
    }

    public function __destruct() {
        Foo::getInstance();
    }
}


$foo = new Foo();
?>
