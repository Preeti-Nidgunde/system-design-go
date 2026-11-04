Design a movie booking system similar to book my show

User - id, name, age
Movie - id, name, category(adult/child)
theater - id, name, loc, address
screen - id, theater_id, name, capacity
seat - id, screen_id, row, no, type
show - id, screen_id, movie_id, start_time, end_time
showPricing - id, show_id, seat_type, price

booking - id, show_id, user_id, total_amount, status
payment - id, booking_id, amount, status (Iniated, refunded, failed)
ShowSeat - id, show_id, seat_id, user_id, booking_id, status unique(show_id, user_id)


Because a booking and a payment have different lifecycles. so keep 2 entities as there will be retries, refunds, reconciliation 

when something can happen multiple times per parent, or has its own lifecycle and external identity, it's an entity, not a column.
