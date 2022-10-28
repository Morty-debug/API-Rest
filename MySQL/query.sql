USE ApiRestDB;

CREATE TABLE tocken(
	idtocken INT NOT NULL AUTO_INCREMENT,
	tocken VARCHAR(200) NOT NULL,
	sistema DATETIME NOT NULL,
	limite DATETIME NOT NULL,
	PRIMARY KEY(idtocken)
);

DELIMITER $$
CREATE FUNCTION creartocken()
    RETURNS VARCHAR(200)
    DETERMINISTIC
    BEGIN
        DECLARE contador INT(11);
        DECLARE token VARCHAR(200);
        DECLARE creacion DATETIME;
        DECLARE temporizado DATETIME;
        SELECT CONCAT(MD5('Tocken'),'.',MD5(RAND()),'.',MD5(NOW())), NOW(), DATE_ADD(NOW(),INTERVAL 10 MINUTE) INTO token, creacion, temporizado;
        SELECT COUNT(*) INTO contador FROM tocken WHERE tocken = token;
        IF contador > 0 THEN
            UPDATE tocken SET sistema=creacion, limite=temporizado WHERE tocken = token;
        ELSE
            INSERT INTO tocken(tocken,sistema,limite) VALUES(token,creacion,temporizado);
        END IF;
        RETURN token;
END$$
DELIMITER ; 


DELIMITER $$
CREATE FUNCTION validartocken(token VARCHAR(200))
    RETURNS VARCHAR(30)
    DETERMINISTIC
    BEGIN
        DECLARE contador INT(11);
        DECLARE limitado INT(11);
        DECLARE retorno VARCHAR(30);
        SELECT COUNT(*), (NOW()>limite) INTO contador, limitado FROM tocken WHERE tocken = token;
        IF contador <= 0 THEN
            SET retorno = 'Token no existe';
        ELSEIF limitado = 1 THEN
            SET retorno = 'Token expirado';
        ELSE
            SET retorno = 'Token Valido';
        END IF;
        RETURN retorno;
END$$
DELIMITER ; 


