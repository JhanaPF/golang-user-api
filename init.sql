-- Database: golang-user-api
CREATE DATABASE golang_user_api
  WITH OWNER = go_user_admin 
  ENCODING = 'UTF8' 
  LIMIT = -1;

-- Table pour les utilisateurs
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    mail VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL
);

-- Table pour les amitiés entre utilisateurs
CREATE TABLE IF NOT EXISTS relationships (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    friend_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (friend_id) REFERENCES users(id)
);

-- Table pour l'historique des leçons terminées
CREATE TABLE IF NOT EXISTS game_history (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(24) NOT NULL,  
    lesson_id VARCHAR(24),         
    course_id VARCHAR(24),  
    game_time INT,  
    score INT,      
    success_rate DECIMAL(5, 2) 
    date DATE
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Table pour les choix de cours
CREATE TABLE IF NOT EXISTS course_choices (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(24) NOT NULL,  
    date DATE
    FOREIGN KEY (user_id) REFERENCES users(id)
);