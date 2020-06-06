DELETE FROM `raw.product` WHERE TRUE;
INSERT INTO `raw.product`
VALUES
    ("100", "product1", "Product 1", "1234567890123", "xyz-1", CURRENT_DATE(), CURRENT_TIMESTAMP()),
    ("200", "product2", "Product 2", "2345678901234", "xyz-2", CURRENT_DATE(), CURRENT_TIMESTAMP())
;

DELETE FROM `raw.basket` WHERE creation_time IS NOT NULL;
INSERT INTO `raw.basket`
VALUES
    ("5000", "12", "3", CURRENT_TIMESTAMP(), FALSE, 4, "EU", 0, 0)
;

DELETE FROM `raw.basket_line` WHERE creation_time IS NOT NULL;
INSERT INTO `raw.basket_line`
VALUES
    ("5000.1", "5000", 1, TIMESTAMP_ADD(CURRENT_TIMESTAMP(), INTERVAL 2 SECOND), FALSE, "100", 1, 13.33, 15.99),
    ("5000.2", "5000", 2, TIMESTAMP_ADD(CURRENT_TIMESTAMP(), INTERVAL 5 SECOND), FALSE, "200", 3, 24.98, 29.97)
;
