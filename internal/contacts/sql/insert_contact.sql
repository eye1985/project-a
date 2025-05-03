insert into contact (user_1, user_2)
VALUES (LEAST($1::int, $2::int),
        GREATEST($1::int, $2::int))
RETURNING id, user_1, user_2;
