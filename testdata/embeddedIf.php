<?php

if (true):
    echo "false";
  else:
    echo "true";
endif;

if (true):
  if (false):
    echo "false";
  else:
    echo "true";
  endif;
endif;
