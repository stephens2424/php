<?php
$tokens = token_get_all('<?php
  class SomeClass {
      const CONST = 1;
      const CONTINUE = (self::CONST + 1);
      const ARRAY = [1, self::CONTINUE => [3, 4], 5];
  }
', TOKEN_PARSE);

array_walk($tokens, function($tk) {
  if(is_array($tk)) {
    if(($t = token_name($tk[0])) == 'T_WHITESPACE') return;
    echo "L{$tk[2]}: ".$t." {$tk[1]}", PHP_EOL;
  }
  else echo $tk, PHP_EOL;
});

echo "Done";

?>
