<?php

class PropelLogger
{

  public function log($message, $severity = null)
  {
    if (true) {
        echo "true";
    }
  }


  function getQueries()
  {
    return $this->queries;
  } 
}
