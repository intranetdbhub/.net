USE `ecsyslabs-network-inventory`;

CREATE TABLE device_inventory (
	`Id` INT AUTO_INCREMENT PRIMARY KEY,
	`Class` VARCHAR(20),
	`Colocation` VARCHAR(30),
	`Description` VARCHAR(200),
	`Hostname` VARCHAR(50),
	`Mgmt IP` VARCHAR(20),
    `Serial #` VARCHAR(30),
    `Type` VARCHAR(30),
    attachment LONGBLOB
);