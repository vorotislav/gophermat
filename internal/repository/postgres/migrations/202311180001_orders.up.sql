CREATE TABLE orders (
    id BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY, -- id записи
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE, -- id пользователя
    order_number TEXT, -- номер заказа
    status TEXT, -- текущий статус заказа
    accrual INT, -- начисления за заказ
    uploaded_at TIMESTAMP WITH TIME ZONE -- отметка времени, когда был загружен заказ
);

CREATE INDEX orders_number ON orders (order_number);