<?php

class myFilterIterator extends FilterIterator {
	function accept() { }
}

class myCachingIterator extends CachingIterator { }

class myRecursiveCachingIterator extends RecursiveCachingIterator { }

class myParentIterator extends ParentIterator { }

class myLimitIterator extends LimitIterator { }

class myNoRewindIterator extends NoRewindIterator  {}

try {
	$it = new myFilterIterator();	
} catch (TypeError $e) {
    echo $e->getMessage(), "\n";
}

try {
	$it = new myCachingIterator();	
} catch (TypeError $e) {
    echo $e->getMessage(), "\n";
}

try {
	$it = new myRecursiveCachingIterator();	
} catch (TypeError $e) {
    echo $e->getMessage(), "\n";
}

try {
	$it = new myParentIterator();	
} catch (TypeError $e) {
    echo $e->getMessage(), "\n";
}

try {
	$it = new myLimitIterator();
} catch (TypeError $e) {
    echo $e->getMessage(), "\n";
}
try {
	$it = new myNoRewindIterator();
} catch (TypeError $e) {
    echo $e->getMessage(), "\n";
}

?>
