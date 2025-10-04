CREATE TABLE Page (
    page_id UUID PRIMARY KEY,
    content TEXT NOT NULL
) PRIMARY KEY (page_id);

CREATE TABLE Paragraph (
    paragraph_id UUID PRIMARY KEY,
    page_id UUID NOT NULL,
    content TEXT NOT NULL,
    bert_vector Array<Float64>,

    CONSTRAINT fk_page FOREIGN KEY (page_id) REFERENCES Page(page_id)
) PRIMARY KEY (paragraph_id);

CREATE OR REPLACE FUNCTION CosineSimilarity(vector1 Array(Float32), vector2 Array(Float32))
RETURNS Float32
AS
$$
WITH dot_product AS (
    SELECT SUM(v1 * v2) as dot
    FROM ARRAY_JOIN(vector1, vector2) AS (v1, v2)
),
norm1 AS (
    SELECT SQRT(SUM(v1 * v1)) as norm
    FROM ARRAY_JOIN(vector1) AS v1
),
norm2 AS (
    SELECT SQRT(SUM(v2 * v2)) as norm
    FROM ARRAY_JOIN(vector2) AS v2
)
SELECT
    CASE
        WHEN norm1.norm = 0 OR norm2.norm = 0 THEN 0.0
        ELSE dot_product.dot / (norm1.norm * norm2.norm)
    END
FROM dot_product, norm1, norm2
$$;
