-- Create/ensure the budgetapp user exists for localhost and any host, set password and grant privileges
CREATE USER IF NOT EXISTS 'budgetapp'@'%' IDENTIFIED BY 'pass12345';
CREATE USER IF NOT EXISTS 'budgetapp'@'localhost' IDENTIFIED BY 'pass12345';
GRANT ALL PRIVILEGES ON budgetapp.* TO 'budgetapp'@'%';
GRANT ALL PRIVILEGES ON budgetapp.* TO 'budgetapp'@'localhost';
FLUSH PRIVILEGES;
