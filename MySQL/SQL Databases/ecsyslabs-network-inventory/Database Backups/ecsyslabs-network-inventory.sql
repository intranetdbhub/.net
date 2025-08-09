CREATE DATABASE  IF NOT EXISTS `ecsyslabs-network-inventory` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `ecsyslabs-network-inventory`;
-- MySQL dump 10.13  Distrib 8.0.43, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: ecsyslabs-network-inventory
-- ------------------------------------------------------
-- Server version	8.0.43

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `device_inventory`
--

DROP TABLE IF EXISTS `device_inventory`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `device_inventory` (
  `Id` int NOT NULL AUTO_INCREMENT,
  `Hostname` varchar(255) DEFAULT NULL,
  `Colocation` varchar(30) DEFAULT NULL,
  `Class` varchar(255) DEFAULT NULL,
  `Description` varchar(200) DEFAULT NULL,
  `Mgmt IP` varchar(20) DEFAULT NULL,
  `Serial #` varchar(30) DEFAULT NULL,
  `Type` varchar(30) DEFAULT NULL,
  `Attachment` longblob,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `device_inventory`
--

LOCK TABLES `device_inventory` WRITE;
/*!40000 ALTER TABLE `device_inventory` DISABLE KEYS */;
INSERT INTO `device_inventory` VALUES (5,'AURCOSTE001A-ERDC001A','AURCOSTE LDC','DELL R-710','ESXi Server','9.9.0.5','N/A','Rack Enterprise Server',NULL),(6,'AURCOSTE001-CSS001A','AURCOSTE LDC','Cisco 3560 Catalyst Switch','LDC Core Switch','10.68.1.2','CAT1123NJ05','Enterprise Switch',NULL),(7,'AURCOSTE001-CSR001A','AURCOSTE LDC','Cisco 2811 Router','LDC Core Router','10.68.2.10','FTX1105A6AZ','Enterprise Router',NULL),(8,'FRECASTE001-STE001A','FRECASTE LDC','GNU/Linux 6.8.0-63-generic x86_64','Ubuntu Cloud Server','23.239.4.97','79840765','Akamai Linode',NULL);
/*!40000 ALTER TABLE `device_inventory` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-08-08 22:36:33
