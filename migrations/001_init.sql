CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    role VARCHAR(20) DEFAULT 'customer',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE bookings (
    id UUID PRIMARY KEY,
    customer_id UUID REFERENCES users(id),
    event_date DATE NOT NULL,
    event_time TIME NOT NULL,
    event_location VARCHAR(500),
    total_price DECIMAL(10,2),
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE payments (
    id UUID PRIMARY KEY,
    booking_id UUID UNIQUE REFERENCES bookings(id),
    amount DECIMAL(10,2),
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE professionals (
                               id UUID PRIMARY KEY,
                               full_name VARCHAR(255) NOT NULL,
                               profession VARCHAR(100),
                               bio TEXT,
                               price_per_hour DECIMAL(10,2),
                               rating_avg DECIMAL(3,2),
                               created_at TIMESTAMP DEFAULT now()
);


CREATE TABLE availability (
                              id UUID PRIMARY KEY,
                              professional_id UUID REFERENCES professionals(id),
                              start_time TIMESTAMP NOT NULL,
                              end_time TIMESTAMP NOT NULL,
                              created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE reviews (
                         id UUID PRIMARY KEY,
                         professional_id UUID REFERENCES professionals(id),
                         customer_id UUID REFERENCES users(id),
                         rating INT CHECK (rating >= 1 AND rating <= 5),
                         comment TEXT,
                         created_at TIMESTAMP DEFAULT now()
);

