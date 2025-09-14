CREATE MATERIALIZED VIEW years_count AS
SELECT 
    date_of_publication_year,
    count(*)
FROM unique_summary_cars
GROUP BY date_of_publication_year
WITH DATA;


CREATE OR REPLACE FUNCTION refresh_years_count()
RETURNS TRIGGER AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY years_count;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_refresh_years_count
AFTER INSERT OR UPDATE OR DELETE ON unique_summary_cars
FOR EACH STATEMENT
EXECUTE FUNCTION refresh_years_count();