<?php
$td = __DIR__ . '/bug66405';
mkdir($td);
touch($td . '/file1.txt');
touch($td . '/file2.md');
mkdir($td . '/testsubdir');
touch($td . '/testsubdir/file3.csv');

class Bug66405 extends RecursiveDirectoryIterator
{
    public function current()
    {
        $current = parent::current();
        echo gettype($current) . " $current\n";
        return $current;
    }

    public function getChildren()
    {
        $children = parent::getChildren();
        if (is_object($children)) {
            echo get_class($children) . " $children\n";
        } else {
            echo gettype($children) . " $children\n";
        }
        return $children;
    }
}

$rdi = new Bug66405($td, FilesystemIterator::CURRENT_AS_PATHNAME | FilesystemIterator::SKIP_DOTS);
$rii = new RecursiveIteratorIterator($rdi);

ob_start();
foreach ($rii as $file) {
    //noop
}
$results = explode("\n", ob_get_clean());
sort($results);
echo implode("\n", $results);
?>
