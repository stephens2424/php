<?php
  class C {}

  function foo($a, $b) {
      echo "Called global foo($a, $b)\n";
  }

  $name = 'functions';

  $c = new C;
  $c->functions[0] = 'foo';
  $c->functions[1][2][3][4] = 'foo';
  
  $c->$name[0](1, 2);
  $c->$name[1][2][3][4](3, 4);
  
  $c->functions[0](5, 6);
  $c->functions[1][2][3][4](7, 8);
?>
