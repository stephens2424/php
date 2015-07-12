<?php

class Date {
    public function __construct($in) {
        $this->date = date_create($in);
    }

    public function getYear1() {
        return date_format($this->date, 'Y');
    }

    public function getYear2() {
        return call_user_func([$this->date, 'format'], 'Y');
    }
}

$d = new Date('NOW');
var_dump($d->getYear1());
var_dump($d->getYear2());

?>
