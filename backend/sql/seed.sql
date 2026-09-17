-- Test User (password: mySuperSecretPassword123)
INSERT INTO users (email, password_hash) 
VALUES ('test@example.com', '$2a$10$ujzR0u6YtpbChC.hrG85s.zdNU1fItcPmoNIXYFraC6gjfgkLeuq.') 
ON CONFLICT (email) DO NOTHING;
