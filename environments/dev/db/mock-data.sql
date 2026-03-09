-- Mock data for development database

INSERT INTO users (email, hashed_password, first_name, last_name) VALUES
  ('admin@example.com', '$2b$12$dummyhashfordev', 'Admin', 'User'),
  ('test@example.com', '$2b$12$dummyhashfordev', 'Test', 'User'),
  ('john.doe@example.com', '$2b$12$dummyhashfordev', 'John', 'Doe'),
  ('jane.smith@example.com', '$2b$12$dummyhashfordev', 'Jane', 'Smith');

INSERT INTO conversations (repo_url, title) VALUES
  ('https://github.com/example/frontend', 'Frontend-projekt'),
  ('https://github.com/example/backend', 'API-utveckling'),
  ('https://github.com/example/devops', 'CI/CD pipeline');

INSERT INTO messages (conversation_id, user_id, message) VALUES
  (1, 1, 'Hej, jag behöver hjälp med att sätta upp Nuxt-projektet'),
  (1, 2, 'Självklart! Vilken version av Node kör du?'),
  (1, 1, 'Jag har Node v20 installerad'),
  (1, 2, 'Perfekt, då kan vi använda Nuxt 4. Låt mig visa dig grundstrukturen.'),
  (1, 1, 'Tack! Vad behöver jag göra först?'),
  (2, 3, 'Kan du hjälpa mig med en API-route?'),
  (2, 2, 'Vad för typ av route behöver du?'),
  (2, 3, 'En GET-route som hämtar användare från databasen'),
  (2, 2, 'Det kan jag hjälpa dig med. Vill du använda SQLAlchemy eller raw SQL?'),
  (3, 4, 'Hej, jag håller på att sätta upp CI/CD'),
  (3, 2, 'Vilken plattform använder du? GitHub Actions?'),
  (3, 4, 'Ja, och jag behöver hjälp med att konfigurera Docker');
