<?php
require_once("connect.inc");
class A {

	private $mysqli;

	public function __construct() {
		global $user, $host, $passwd, $db, $port, $socket;
		$this->mysqli = new mysqli($host, $user, $passwd, $db, $port, $socket);
		$result = $this->mysqli->query("SELECT NOW() AS my_time FROM DUAL");
		$row = $result->fetch_object();
		echo $row->my_time."<br>\n";
		$result->close();
	}

	public function __destruct() {
		$this->mysqli->close();
	}
}

class B {

	private $mysqli;

	public function __construct() {
		global $user, $host, $passwd, $db, $port, $socket;
		$this->mysqli = new mysqli($host, $user, $passwd, $db, $port, $socket);
		$result = $this->mysqli->query("SELECT NOW() AS my_time FROM DUAL");
		$row = $result->fetch_object();
		echo $row->my_time."<br>\n";
		$result->close();
	}

	public function __destruct() {
		$this->mysqli->close();
	}
}

$A = new A();
$B = new B();
?>
