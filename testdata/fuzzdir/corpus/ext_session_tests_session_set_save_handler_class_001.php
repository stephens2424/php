<?php

ob_start();

/* 
 * Prototype : bool session_set_save_handler(SessionHandler $handler [, bool $register_shutdown_function = true])
 * Description : Sets user-level session storage functions
 * Source code : ext/session/session.c 
 */

echo "*** Testing session_set_save_handler() : basic class wrapping existing handler ***\n";

class MySession extends SessionHandler {
	public $i = 0;
	public function open($path, $name) {
		++$this->i;
		echo 'Open ', session_id(), "\n";
		return parent::open($path, $name);
	}
	public function create_sid() {
		// This method should be removed when 5.5 become unsupported.
		++$this->i;
		echo 'Old Create SID ', session_id(), "\n";
		return parent::create_sid();
	}
	public function read($key) {
		++$this->i;
		echo 'Read ', session_id(), "\n";
		return parent::read($key);
	}
	public function write($key, $data) {
		++$this->i;
		echo 'Write ', session_id(), "\n";
		return parent::write($key, $data);
	}
	public function close() {
		++$this->i;
		echo 'Close ', session_id(), "\n";
		return parent::close();
	}
	public function createSid() {
		// User should use this rather than create_sid()
		// If both create_sid() and createSid() exists,
		// createSid() is used.
		++$this->i;
		echo 'New Create ID ', session_id(), "\n";
		return parent::create_sid();
	}
	public function validateId($key) {
		++$this->i;
		echo 'Validate ID ', session_id(), "\n";
		return TRUE;
		// User must implement their own method and
		// cannot call parent as follows.
		// return parent::validateSid($key);
	}
	public function updateTimestamp($key, $data) {
		++$this->i;
		echo 'Update Timestamp ', session_id(), "\n";
		return parent::write($key, $data);
		// User must implement their own method and
		// cannot call parent as follows
		// return parent::updateTimestamp($key, $data);
	}
}

$oldHandler = ini_get('session.save_handler');
$handler = new MySession;
session_set_save_handler($handler);
session_start();

var_dump(session_id(), $oldHandler, ini_get('session.save_handler'), $handler->i, $_SESSION);

$_SESSION['foo'] = "hello";

session_write_close();
session_unset();

session_start();
var_dump($_SESSION);

session_write_close();
session_unset();
var_dump($handler->i);

