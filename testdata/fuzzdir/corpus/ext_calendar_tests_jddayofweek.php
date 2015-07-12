<?php
foreach (array(2440588, 2452162, 2453926, -1000) as $jd) {
  echo "### JD $jd ###\n";
  foreach (array(CAL_DOW_DAYNO, CAL_DOW_LONG, CAL_DOW_SHORT) as $mode) {
    echo "--- mode $mode ---\n";
    for ($offset = 0; $offset <= 7; $offset++) {
      echo jddayofweek($jd + $offset, $mode). "\n";
    }
  }
}
?>
