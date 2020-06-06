#standardSQL
WITH
    lines_by_basket_id AS (
        SELECT
            basket_id AS id,
            ARRAY_AGG(
                STRUCT(
                    id AS basket_line_id,
                    basket_line,
                    creation_time,
                    cancelled,
                    product_id,
                    nb_items,
                    net_value_wo_tax,
                    net_value_w_tax
                )
                ORDER BY basket_line
            ) AS lines
        FROM
            `raw.basket_line`
        WHERE
            DATE(creation_time) = CURRENT_DATE()
        GROUP BY
            basket_id
    )
SELECT
    *
FROM
    `raw.basket`
LEFT JOIN lines_by_basket_id USING (id)
WHERE
    DATE(creation_time) = CURRENT_DATE()
