<?php
echo "Hello world\n";
var_dump($argv);

echo "---\n";

while($f = fgets(STDIN)){
    echo "stdin: $f";
}
