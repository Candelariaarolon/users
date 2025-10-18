-- MySQL dump 10.13  Distrib 8.0.42, for Win64 (x86_64)
--
-- Host: localhost    Database: arqui_software
-- ------------------------------------------------------
-- Server version	9.3.0

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
-- Table structure for table `activity_models`
--

DROP TABLE IF EXISTS `activity_models`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `activity_models` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `category` varchar(100) NOT NULL,
  `name` varchar(100) NOT NULL,
  `description` varchar(350) NOT NULL,
  `profesor_name` varchar(100) NOT NULL,
  `quotas` bigint NOT NULL,
  `day` varchar(50) NOT NULL,
  `hour_start` varchar(50) DEFAULT NULL,
  `active` tinyint(1) DEFAULT '1',
  `photo` varchar(300) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `activity_models`
--

LOCK TABLES `activity_models` WRITE;
/*!40000 ALTER TABLE `activity_models` DISABLE KEYS */;
INSERT INTO `activity_models` VALUES (1,'Deportes','Tenis','Clases de tenis','María López',15,'Jueves','17:00:00',1,'https://tse4.mm.bing.net/th/id/OIP.AEYs1txwFIeMfu2VDCaTuQHaEK?rs=1&pid=ImgDetMain&o=7&rm=3'),(2,'Box','Musculacion','Gimnasio box libre','Carlos Gómez',10,'Miércoles','19:00:00',1,'https://tse3.mm.bing.net/th/id/OIP.JAnfGoGxFFPacYpz6GHdrgHaEK?rs=1&pid=ImgDetMain&o=7&rm=3'),(3,'Cardio','Spinning','Clases de spinning sobre bicicleta fija y musica','Lucía Martínez',25,'Lunes','16:00:00',1,'https://tse3.mm.bing.net/th/id/OIP.3wmV06Ok8naChVrk0W8IlgHaEK?rs=1&pid=ImgDetMain&o=7&rm=3'),(4,'Salud','Yoga y meditación','Sesiones semanales para mejorar flexibilidad y reducir el estrés.','Ana Fernández',12,'Viernes','20:00:00',1,'https://tse2.mm.bing.net/th/id/OIP.N5FcRhhNbKfR1T-d3SfwzQHaE8?rs=1&pid=ImgDetMain&o=7&rm=3'),(5,'Cardio','GAP','Gluteos, Abdomen y Piernas','Geronimo Benavidez',20,'Martes','10:00:00',1,'https://medac.es/sites/default/files/blog/destacadas/Conoce%206%20ejercicios%20de%20GAP%20para%20hacer%20en%2030%20minutos-min.jpg');
/*!40000 ALTER TABLE `activity_models` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-06-29 17:55:49
